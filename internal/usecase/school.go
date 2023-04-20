package usecase

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/googlecloud"
)

type JobProgress struct {
	Epoch       int
	IsCompleted bool
}

type AIJobInstance struct {
	u            Usecase
	job          *entity.AISchoolJob
	currentEpoch int
	IsCompleted  bool
	progCh       chan JobProgress
}

var currentAIJobs map[string]AIJobInstance

func (u Usecase) JobAIS_WatchPending() error {
	jobList, err := u.Repo.GetAISchoolJobByStatus([]string{"running", "waiting"})
	if err != nil {
		return err
	}

	if currentAIJobs == nil {
		currentAIJobs = make(map[string]AIJobInstance)
	}
	for jobID, job := range currentAIJobs {
		if job.IsCompleted {
			delete(currentAIJobs, jobID)
		}
	}
	if len(currentAIJobs) >= 10 {
		return nil
	}
	for _, job := range jobList {
		if job.Status == "waiting" {
			jobParams := &structure.AISchoolModelParams{}
			err := json.Unmarshal([]byte(job.Params), jobParams)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while unmarshalling job params: "+err.Error(), "error")
				continue
			}

			newJob := AIJobInstance{
				u:   u,
				job: &job,
			}
			currentAIJobs[job.JobID] = newJob
			job.Status = "running"
			job.ExecutedAt = time.Now().Unix()
			err = u.Repo.UpdateAISchoolJob(&job)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
				continue
			}
			go newJob.Start()
		}
		if job.Status == "running" {
			if _, exist := currentAIJobs[job.JobID]; !exist {
				job.Status = "waiting"
				job.ExecutedAt = 0
				job.Progress = 0
				err = u.Repo.UpdateAISchoolJob(&job)
				if err != nil {
					// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
					continue
				}
			}
		}
	}
	return nil
}

const basePath = "./ai-school-work/"

func createAISchoolWorkFolder(jobID string, params structure.AISchoolModelParams) error {
	if err := os.MkdirAll(basePath+jobID, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(basePath+jobID+"/dataset", os.ModePerm); err != nil {
		return err
	}

	return nil
}

func prepAISchoolWorkFolder(jobID string, params structure.AISchoolModelParams, datasetGCPath string, gcs googlecloud.IGcstorage) error {
	err := createAISchoolWorkFolder(jobID, params)
	if err != nil {
		return err
	}

	content, err := json.Marshal(params)
	if err != nil {
		return err
	}
	log.Println("Writing params to file: ", basePath+jobID+"/params.json")
	err = ioutil.WriteFile(basePath+jobID+"/params.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Unzipping dataset: ", datasetGCPath)
	dataseBytes, err := gcs.ReadFileFromBucketAbs(datasetGCPath)
	if err != nil {
		return err
	}
	log.Println("Dataset size: ", len(dataseBytes))
	br := bytes.NewReader(dataseBytes)

	zr, err := zip.NewReader(br, int64(len(dataseBytes)))
	if err != nil {
		return err
	}
	destination, err := filepath.Abs(basePath + jobID + "/dataset")
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		log.Println("Unzipping", f.Name)
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func clearAISchoolWorkFolder(jobID string) error {
	err := os.RemoveAll("./ai-school-work/" + jobID + "/")
	if err != nil {
		return err
	}
	return nil
}

func (job *AIJobInstance) Start() {
	progCh := make(chan JobProgress)
	jobID := job.job.JobID
	defer func() {
		job.IsCompleted = true
		close(progCh)
		clearAISchoolWorkFolder(jobID)
	}()
	log.Println("Starting job: ", jobID)
	err := clearAISchoolWorkFolder(jobID)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}
	params := structure.AISchoolModelParams{}
	err = json.Unmarshal([]byte(job.job.Params), &params)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}

	dataset, err := job.u.Repo.GetFileByUUID(job.job.DatasetUUID)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}

	err = prepAISchoolWorkFolder(jobID, params, dataset.FileName, job.u.GCS)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}

	job.progCh = progCh
	go func() {
		for prog := range job.progCh {
			job.job.Progress = prog.Epoch
			err = job.u.Repo.UpdateAISchoolJob(job.job)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
				log.Println(err)
			}
		}
	}()
	scriptPath := os.Getenv("AI_SCHOOL_SCRIPT")
	jobPath := basePath + jobID
	jobPathAbs, err := filepath.Abs(jobPath)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}

	jobLog, jobErrLog, err := executeAISchoolJob(scriptPath, jobPathAbs+"/params.json", jobPathAbs+"/dataset", jobPathAbs+"/output.json", job.progCh)
	job.job.Logs = jobLog
	job.job.ErrLogs = jobErrLog
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}
	cloudPath := fmt.Sprintf("ai-school/%s", job.job.JobID)
	uploaded, err := job.u.GCS.FileUploadToBucketInternal(jobPathAbs+"/output.json", &cloudPath)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}

	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	fileModel := &entity.Files{
		FileName: uploaded.Name,
		FileSize: int(uploaded.Size),
		MineType: uploaded.Minetype,
		URL:      cdnURL,
	}

	err = job.u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		job.job.Errors = err.Error()
		job.job.Status = "error"
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
		return
	}
	job.job.OutputUUID = fileModel.UUID
	job.job.OutputLink = fileModel.URL
	job.job.CompletedAt = time.Now().Unix()
	job.job.Status = "completed"
	err = job.u.Repo.UpdateAISchoolJob(job.job)
	if err != nil {
		// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
		return
	}
}
func executeAISchoolJob(scriptPath string, params string, dataset string, output string, progCh chan JobProgress) (string, string, error) {
	// 1. Get params
	// 2. Get dataset
	// 3. Run job
	// 4. Update job
	jobLog := ""
	jobErrLog := ""
	args := fmt.Sprintf("%v -c %v -d %v -o %v", scriptPath, params, dataset, output)
	cmd := exec.Command("python3", strings.Split(args, " ")...)
	// cmd := exec.Command("ls", "-a")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return jobLog, jobErrLog, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return jobLog, jobErrLog, err
	}
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	errStr := ""
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println("err", m)
		jobErrLog += fmt.Sprintln(m)
	}

	scanner2 := bufio.NewScanner(stdout)
	scanner2.Split(bufio.ScanLines)
	for scanner2.Scan() {
		m := scanner2.Text()
		jobLog += fmt.Sprintln(m)
		if strings.Contains(strings.ToLower(m), "epoch") {
			epochStr := strings.Split(m, "Epoch ")
			epochs := strings.Split(epochStr[1], "/")
			currentEpoch := epochs[0]
			currentEpochInt, err := strconv.ParseInt(currentEpoch, 10, 64)
			if err != nil {
				errStr += fmt.Sprintln(err.Error())
				continue
			}
			progCh <- JobProgress{
				Epoch: int(currentEpochInt),
			}
		}
	}

	cmd.Wait()
	if len(errStr) > 0 {
		return jobLog, jobErrLog, errors.New(errStr)
	}
	time.Sleep(100 * time.Millisecond)
	return jobLog, jobErrLog, nil
}