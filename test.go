package main

import (
	"fmt"
	"encoding/json"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type ProductParameters struct {
	ID 			string	`json:"ID"`
	Temperature int 	`json:"Temperature"`
	Humidity 	int 	`json:"Humidity"`
	Vibration 	int 	`json:"Vibration"`
	Shock 		bool 	`json:"Shock"`
	Location 	int 	`json:"Location"`
}

//Initialize the ledger, Genesis Asset
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []ProductParameters{
		{ID: "shipment1", Temperature: 74, Humidity: 85, Vibration: 0, Shock: false, Location: 22.36254},
	}
	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
	err = ctx.GetStub().PutState(asset.ID, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	}
	return nil
}

//Create Asset
func (s *SmartContract) CreateAsset(ctx TransactionContextInterface, id string, temperature int, humidity int, vibration int, shock bool, location int) error {
	exists, err := s.AssetExists(ctx,id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("te asset %s already exists",id)
	}
	asset := ProductParameters{
		ID:				id,
		Temperature: 	temperature,
		Humidity:		humidity,
		Vibration:		vibration,
		Shock:			shock,
		Location:		location,
	}
	assetJSON, err := json.Marshal(asset)
	if err!= nil {
		return err
	}
	return ctx.GetStub().PutState(id, assetJSON)
}

//Query Asset
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*ProductParameters, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("The asset %s does not exist", id)
	}
	var asset ProductParameters
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

//Update Asset
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, temperature int, humidity int, vibration int, shock bool, location int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	asset := ProductParameters{
		ID:				id,
		Temperature: 	temperature,
		Humidity: 		humidity,
		Vibration: 		vibration,
		Shock: 			shock,
		Location: 		location,
	}
	assetJSON, err := json.Marshall(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, assetJSON)
}

//Delete Asset
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

//Check is Asset Exists
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

//NOT COMPLETE ADD OWNER
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx,id)
	if err != nil {
		return err
	}
	asset.Owner = newOwner
	assetJSON ,err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, assetJSON)
}

//Query for all the assets on the ledger
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*ProductParameters, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*ProductParameters
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset ProductParameters
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		asset = append(assets, &asset)
	}
	return assets, nil
}

