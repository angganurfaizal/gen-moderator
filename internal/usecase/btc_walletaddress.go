package usecase

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/url"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/jinzhu/copier"

// 	"rederinghub.io/external/ord_service"
// 	"rederinghub.io/internal/entity"
// 	"rederinghub.io/internal/usecase/structure"
// 	"rederinghub.io/utils"
// 	"rederinghub.io/utils/btc"
// 	"rederinghub.io/utils/helpers"
// )

// func (u Usecase) CreateOrdBTCWalletAddress(input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))

// 	// find Project and make sure index < max supply
// 	projectID := input.ProjectID
// 	project, err := u.Repo.FindProjectByProjectIdWithoutCache(projectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	if project.MintingInfo.Index >= project.MaxSupply {
// 		err = fmt.Errorf("project %s is minted out", projectID)
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	walletAddress := &entity.BTCWalletAddress{}
// 	err = copier.Copy(walletAddress, input)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)

// 	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
// 		Args: []string{
// 			"--wallet",
// 			userWallet,
// 			"wallet",
// 			"create",
// 		},
// 	})

// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		//return nil, err
// 	} else {
// 		walletAddress.Mnemonic = resp.Stdout
// 	}

// 	logger.AtLog.Logger.Info("CreateOrdBTCWalletAddress.createdWallet", zap.Any("resp", resp))
// 	resp, err = u.OrdService.Exec(ord_service.ExecRequest{
// 		Args: []string{
// 			"--wallet",
// 			userWallet,
// 			"wallet",
// 			"receive",
// 		},
// 	})
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("CreateOrdBTCWalletAddress.receive", zap.Any("resp", resp))
// 	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("found.Project", zap.Any("p.ID", p.ID))
// 	mintPrice, err := strconv.Atoi(p.MintPrice)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	networkFee, err := strconv.Atoi(p.NetworkFee)
// 	if err == nil {
// 		mintPrice += networkFee
// 	}
// 	walletAddress.Amount = strconv.Itoa(mintPrice)
// 	walletAddress.UserAddress = userWallet
// 	walletAddress.OriginUserAddress = input.WalletAddress
// 	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
// 	walletAddress.IsConfirm = false
// 	walletAddress.IsMinted = false
// 	walletAddress.FileURI = ""       //find files from google store
// 	walletAddress.InscriptionID = "" //find files from google store
// 	walletAddress.ProjectID = input.ProjectID
// 	walletAddress.Balance = "0"
// 	walletAddress.BalanceCheckTime = 0

// 	err = u.Repo.InsertBtcWalletAddress(walletAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	return walletAddress, nil
// }

// func (u Usecase) CreateSegwitBTCWalletAddress(input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
// 	walletAddress := &entity.BTCWalletAddress{}
// 	privKey, _, addressSegwit, err := btc.GenerateAddressSegwit()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	walletAddress.OrdAddress = addressSegwit //TODO: @thaibao/@tri check this field
// 	walletAddress.Mnemonic = privKey
// 	walletAddress.UserAddress = helpers.CreateBTCOrdWallet(input.WalletAddress)
// 	logger.AtLog.Logger.Info("CreateSegwitBTCWalletAddress.receive", zap.Any("addressSegwit", addressSegwit))
// 	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("found.Project", zap.Any("p.ID", p.ID))
// 	mintPrice, err := strconv.Atoi(p.MintPrice)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	networkFee, err := strconv.Atoi(p.NetworkFee)
// 	if err == nil {
// 		mintPrice += networkFee
// 	}

// 	walletAddress.Amount = strconv.Itoa(mintPrice)
// 	walletAddress.OriginUserAddress = input.WalletAddress
// 	walletAddress.IsConfirm = false
// 	walletAddress.IsMinted = false
// 	walletAddress.FileURI = ""       //find files from google store
// 	walletAddress.InscriptionID = "" //find files from google store
// 	walletAddress.ProjectID = input.ProjectID
// 	walletAddress.Balance = "0"
// 	walletAddress.BalanceCheckTime = 0

// 	err = u.Repo.InsertBtcWalletAddress(walletAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	return walletAddress, nil
// }

// func (u Usecase) CheckBalanceWalletAddress(input structure.CheckBalance) (*entity.BTCWalletAddress, error) {

// 	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	balance, err := u.CheckBalance(*btc)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	return balance, nil
// }

// func (u Usecase) BTCMint(input structure.BctMintData) (*ord_service.MintStdoputRespose, *string, error) {
// 	eth := &entity.ETHWalletAddress{}
// 	mintype := entity.BIT
// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))

// 	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
// 	if err != nil {
// 		btc = &entity.BTCWalletAddress{}
// 		eth, err = u.Repo.FindEthWalletAddressByOrd(input.Address)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}

// 		err = copier.Copy(btc, eth)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}

// 		mintype = entity.ETH
// 	}

// 	btc, err = u.MintLogic(btc)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}

// 	// get data from project
// 	p, err := u.Repo.FindProjectByTokenID(btc.ProjectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}
// 	//logger.AtLog.Logger.Info("found.Project", zap.Any("p", p))

// 	//prepare data for mint
// 	// - Get project.AnimationURL
// 	projectNftTokenUri := &structure.ProjectAnimationUrl{}
// 	err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}

// 	// - Upload the Animation URL to GCS
// 	animation := projectNftTokenUri.AnimationUrl
// 	logger.AtLog.Logger.Info("animation", zap.Any("animation", animation))
// 	if animation != "" {
// 		animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
// 		now := time.Now().UTC().Unix()
// 		uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", p.TokenID, now))
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}
// 		btc.FileURI = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)

// 	} else {
// 		images := p.Images
// 		logger.AtLog.Logger.Info("images", zap.Any("len(images)", len(images)))
// 		if len(images) > 0 {
// 			btc.FileURI = images[0]
// 			newImages := []string{}
// 			processingImages := p.ProcessingImages
// 			//remove the project's image out of the current projects
// 			for i := 1; i < len(images); i++ {
// 				newImages = append(newImages, images[i])
// 			}
// 			processingImages = append(p.ProcessingImages, btc.FileURI)
// 			p.Images = newImages
// 			p.ProcessingImages = processingImages
// 		}
// 	}
// 	//end Animation URL
// 	if btc.FileURI == "" {
// 		err = errors.New("There is no file uri to mint")
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}

// 	baseUrl, err := url.Parse(btc.FileURI)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}

// 	mintData := ord_service.MintRequest{
// 		WalletName: os.Getenv("ORD_MASTER_ADDRESS"),
// 		FileUrl:    baseUrl.String(),
// 		FeeRate:    entity.DEFAULT_FEE_RATE, //temp
// 		DryRun:     false,
// 	}

// 	logger.AtLog.Logger.Info("mintData", zap.Any("mintData", mintData))
// 	resp, err := u.OrdService.Mint(mintData)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}
// 	logger.AtLog.Logger.Info("mint.resp", zap.Any("resp", resp))
// 	//update btc or eth here
// 	if mintype == entity.BIT {
// 		btc.IsMinted = true
// 		btc.FileURI = baseUrl.String()
// 		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, btc)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}
// 		logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))

// 	} else {
// 		eth.IsMinted = true
// 		eth.FileURI = baseUrl.String()
// 		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(eth.OrdAddress, eth)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}
// 		logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 	}

// 	updated, err := u.Repo.UpdateProject(p.UUID, p)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}
// 	logger.AtLog.Logger.Info("project.Updated", zap.Any("updated", updated))

// 	u.Notify(fmt.Sprintf("[MintFor][%s][projectID %s]", mintype, btc.ProjectID), btc.OrdAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", btc.UserAddress, resp.Stdout))
// 	tmpText := resp.Stdout
// 	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
// 	jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
// 	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
// 	btcMintResp := &ord_service.MintStdoputRespose{}

// 	bytes := []byte(jsonStr)
// 	err = json.Unmarshal(bytes, btcMintResp)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, nil, err
// 	}

// 	if mintype == entity.BIT {
// 		btc.MintResponse = entity.MintStdoputResponse(*btcMintResp)
// 		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, btc)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}
// 		logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))

// 	} else {
// 		eth.MintResponse = entity.MintStdoputResponse(*btcMintResp)
// 		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(eth.OrdAddress, eth)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			return nil, nil, err
// 		}
// 		logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 	}

// 	u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
// 		TokenID:       btcMintResp.Inscription,
// 		Commit:        btcMintResp.Commit,
// 		Reveal:        btcMintResp.Reveal,
// 		Fees:          btcMintResp.Fees,
// 		MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
// 		Owner:         "",
// 		ProjectID:     btc.ProjectID,
// 		Action:        entity.MINT,
// 		Type:          mintype,
// 		Balance:       btc.Balance,
// 		Amount:        btc.Amount,
// 		ProccessID:    btc.UUID,
// 	})

// 	return btcMintResp, &btc.FileURI, nil
// }

// func (u Usecase) ReadGCSFolder(input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))
// 	u.GCS.ReadFolder("btc-projects/project-1/")
// 	return nil, nil
// }

// func (u Usecase) UpdateBtcMintedStatus(btcWallet *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
// 	logger.AtLog.Logger.Info("input", zap.Any("btcWallet", btcWallet))

// 	btcWallet.IsMinted = true

// 	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btcWallet.OrdAddress, btcWallet)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 	return btcWallet, nil
// }

// func (u Usecase) GetBalanceSegwitBTCWallet(userAddress string) (string, error) {

// 	logger.AtLog.Logger.Info("userAddress", zap.Any("userAddress", userAddress))

// 	_, bs, err := u.buildBTCClient()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return "", nil
// 	}
// 	logger.AtLog.Logger.Info("bs", zap.Any("bs", bs))
// 	balance, confirm, err := bs.GetBalance(userAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return "", err
// 	}
// 	logger.AtLog.Logger.Info("confirm", zap.Any("confirm", confirm))
// 	logger.AtLog.Logger.Info("balance", zap.Any("balance.String()", balance.String()))

// 	//TODO: @thaibao

// 	_ = confirm

// 	return balance.String(), nil
// }

// func (u Usecase) GetBalanceOrdBTCWallet(userAddress string) (string, error) {
// 	balanceRequest := ord_service.ExecRequest{
// 		Args: []string{
// 			"--wallet",
// 			userAddress,
// 			"wallet",
// 			"balance",
// 		},
// 	}

// 	logger.AtLog.Logger.Info("balanceRequest", zap.Any("balanceRequest", balanceRequest))
// 	//userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
// 	resp, err := u.OrdService.Exec(balanceRequest)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return "", err
// 	}

// 	logger.AtLog.Logger.Info("balanceResponse", zap.Any("resp", resp))
// 	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
// 	return balance, nil
// }

// func (u Usecase) CheckBalance(btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {

// 	//TODO - removed checking ORD, only Segwit is used
// 	balance, err := u.GetBalanceSegwitBTCWallet(btc.OrdAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	if balance == "" {
// 		err := errors.New("balance is empty")
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("balance", zap.Any("balance", balance))
// 	btc.Balance = strings.ReplaceAll(balance, `\n`, "")
// 	btc.BalanceCheckTime = btc.BalanceCheckTime + 1
// 	updated, _ := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, &btc)
// 	logger.AtLog.Logger.Info("updated", zap.Any("btc", btc))
// 	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 	return &btc, nil
// }

// func (u Usecase) BalanceLogic(btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
// 	balance, err := u.CheckBalance(btc)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	//TODO logic of the checked balance here
// 	if balance.Balance < btc.Amount {
// 		err := errors.New("Not enough amount")
// 		return nil, err
// 	}
// 	btc.IsConfirm = true

// 	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, &btc)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 	return &btc, nil
// }

// func (u Usecase) MintLogic(btc *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
// 	var err error

// 	//if this was minted, skip it
// 	if btc.IsMinted {
// 		err = errors.New("This btc was minted")
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	if !btc.IsConfirm {
// 		err = errors.New("This btc must be IsConfirmed")
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	if btc.MintResponse.Inscription != "" {
// 		err = errors.New(fmt.Sprintf("This btc has Inscription %s", btc.MintResponse.Inscription))
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}
// 	if btc.MintResponse.Reveal != "" {
// 		err = errors.New(fmt.Sprintf("This btc has revealID %s", btc.MintResponse.Reveal))
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("btc", zap.Any("btc", btc))
// 	return btc, nil
// }

// // Mint flow
// func (u Usecase) WaitingForBalancing() ([]entity.BTCWalletAddress, error) {
// 	addreses, err := u.Repo.ListProcessingWalletAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {
// 		func(item entity.BTCWalletAddress) {

// 			newItem, err := u.BalanceLogic(item)
// 			if err != nil {
// 				logger.AtLog.Logger.Error("err", zap.Error(err))
// 				return
// 			}
// 			logger.AtLog.Logger.Info(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)
// 			u.Notify(fmt.Sprintf("[WaitingForBalancing][projectID %s]", item.ProjectID), item.UserAddress, fmt.Sprintf("%s checkint BTC %s from [user_address] %s", item.OrdAddress, item.Balance, item.UserAddress))

// 			updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, newItem)
// 			if err != nil {
// 				logger.AtLog.Logger.Error("err", zap.Error(err))
// 				return
// 			}
// 			logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
// 				MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
// 				Owner:         "",
// 				ProjectID:     item.ProjectID,
// 				Action:        entity.BLANCE,
// 				Type:          entity.BIT,
// 				Balance:       item.Balance,
// 				Amount:        item.Amount,
// 				ProccessID:    item.UUID,
// 			})
// 		}(item)

// 		time.Sleep(2 * time.Second)
// 	}

// 	return nil, nil
// }

// func (u Usecase) WaitingForMinting() ([]entity.BTCWalletAddress, error) {
// 	addreses, err := u.Repo.ListMintingWalletAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {
// 		func(item entity.BTCWalletAddress) {

// 			if item.MintResponse.Inscription != "" {
// 				err = errors.New("Token is being minted")
// 				logger.AtLog.Logger.Error("err", zap.Error(err))
// 				return
// 			}

// 			_, _, err := u.BTCMint(structure.BctMintData{Address: item.OrdAddress})
// 			if err != nil {
// 				u.Notify(fmt.Sprintf("[Error][MintFor][projectID %s]", item.ProjectID), item.OrdAddress, err.Error())
// 				logger.AtLog.Logger.Error("err", zap.Error(err))
// 				return
// 			}

// 		}(item)

// 		time.Sleep(2 * time.Second)
// 	}

// 	return nil, nil
// }

// func (u Usecase) WaitingForMinted() ([]entity.BTCWalletAddress, error) {

// 	_, bs, err := u.buildBTCClient()

// 	addreses, err := u.Repo.ListBTCAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {
// 		func(item entity.BTCWalletAddress) {

// 			addr := item.OriginUserAddress
// 			if addr == "" {
// 				addr = item.UserAddress
// 			}

// 			//check token is created or not via BlockcypherService
// 			txInfo, err := bs.CheckTx(item.MintResponse.Reveal)
// 			if err != nil {
// 				logger.AtLog.Logger.Error("err", zap.Error(err))
// 				u.Notify(fmt.Sprintf("[Error][BTC][SendToken.bs.CheckTx][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %s", err.Error(), item.UUID))
// 				return
// 			}
// 			logger.AtLog.Logger.Info("txInfo", zap.Any("txInfo", txInfo))
// 			if txInfo.Confirmations > 1 {
// 				sentTokenResp, err := u.SendToken(addr, item.MintResponse.Inscription)
// 				if err != nil {
// 					u.Notify(fmt.Sprintf("[Error][BTC][SendToken][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %s", err.Error(), item.UUID))
// 					logger.AtLog.Logger.Error("err", zap.Error(err))
// 					return
// 				}

// 				logger.AtLog.Logger.Info(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)

// 				u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
// 					TokenID:       item.MintResponse.Inscription,
// 					Commit:        item.MintResponse.Commit,
// 					Reveal:        item.MintResponse.Reveal,
// 					Fees:          item.MintResponse.Fees,
// 					MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
// 					Owner:         item.UserAddress,
// 					Action:        entity.SENT,
// 					ProjectID:     item.ProjectID,
// 					Type:          entity.BIT,
// 					Balance:       item.Balance,
// 					Amount:        item.Amount,
// 					ProccessID:    item.UUID,
// 				})

// 				u.Notify(fmt.Sprintf("[SendToken][ProjectID: %s]", item.ProjectID), addr, item.MintResponse.Inscription)

// 				// logger.AtLog.Logger.Info("fundResp", fundResp
// 				item.MintResponse.IsSent = true
// 				updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, &item)
// 				if err != nil {
// 					logger.AtLog.Logger.Error("err", zap.Error(err))
// 					return
// 				}

// 				//TODO: - create entity.TokenURI
// 				_, err = u.CreateBTCTokenURI(item.ProjectID, item.MintResponse.Inscription, item.FileURI, entity.BIT)
// 				if err != nil {
// 					logger.AtLog.Logger.Error("err", zap.Error(err))
// 					return
// 				}
// 				logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 				err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.MintResponse.Inscription)
// 				if err != nil {
// 					logger.AtLog.Logger.Error("err", zap.Error(err))
// 					return
// 				}
// 				go u.CreateMintActivity(item.InscriptionID, item.Amount)
// 				go u.NotifyNFTMinted(item.OriginUserAddress, item.InscriptionID, item.MintResponse.Fees)

// 			} else {
// 				logger.AtLog.Logger.Info("checkTx.Inscription.Existed", zap.Any("false", false))
// 			}
// 		}(item)

// 		time.Sleep(5 * time.Second)
// 	}

// 	return nil, nil
// }

// //End Mint flow

// func (u Usecase) SendToken(receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {

// 	sendTokenReq := ord_service.ExecRequest{
// 		Args: []string{
// 			"--wallet",
// 			os.Getenv("ORD_MASTER_ADDRESS"),
// 			"wallet",
// 			"send",
// 			receiveAddr,
// 			inscriptionID,
// 			"--fee-rate",
// 			fmt.Sprintf("%d", entity.DEFAULT_FEE_RATE),
// 		}}

// 	logger.AtLog.Logger.Info("sendTokenReq", zap.Any("sendTokenReq", sendTokenReq))
// 	resp, err := u.OrdService.Exec(sendTokenReq)

// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("sendTokenRes", zap.Any("resp", resp))
// 	return resp, err

// }

// func (u Usecase) Notify(title string, userAddress string, content string) {

// 	//slack
// 	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), "", userAddress)
// 	c := fmt.Sprintf("%s", content)

// 	if _, _, err := u.Slack.SendMessageToSlack(preText, title, c); err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 	}
// }

// func (u Usecase) NotifyWithChannel(channelID string, title string, userAddress string, content string) {
// 	//slack
// 	preText := fmt.Sprintf("[App: %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), userAddress)
// 	c := fmt.Sprintf("%s", content)

// 	if _, _, err := u.Slack.SendMessageToSlackWithChannel(channelID, preText, title, c); err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 	}
// }

// // phuong:
// // send btc from segwit address to master address - it does not call our ORD server
// func (u Usecase) JobBtcSendBtcToMaster() error {

// 	addreses, err := u.Repo.ListWalletAddressToClaimBTC()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("err", zap.Error(err))
// 		return err
// 	}
// 	_, bs, err := u.buildBTCClient()

// 	if err != nil {
// 		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
// 		return err
// 	}

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {

// 		// send master now:
// 		tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(item.Mnemonic, item.OrdAddress, utils.MASTER_ADDRESS, -1, btc.PreferenceMedium)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			continue
// 		}
// 		// save tx:
// 		item.TxSendMaster = tx
// 		item.IsSentMaster = true
// 		_, err = u.Repo.UpdateBtcWalletAddress(&item)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("err", zap.Error(err))
// 			continue
// 		}

// 		time.Sleep(3 * time.Second)
// 	}

// 	return nil
// }

