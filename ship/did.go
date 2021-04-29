package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Iotic-Labs/iotic-sdk-identity/sdk/go/identity"
	log "github.com/sirupsen/logrus"
)

type Did struct {
	masterBytes  []byte
	agentKeyName string
	agentDid     string
	agentPurpose string
	twinPurpose  string
}

func NewDidFromConfig(config Config) (*Did, error) {
	return NewDid(config.AgentSeed, config.AgentKeyName, config.AgentDid)
}

func NewDid(seedString string, agentKeyName string, agentDid string) (*Did, error) {
	method := identity.SeedMethodBip39
	//method := identity.SeedMethodNone

	seed, err := hex.DecodeString(seedString)
	if err != nil {
		return nil, err
	}

	masterBytes, err := identity.SeedToMaster(seed, "", method)
	if err != nil {
		return nil, err
	}
	rtn := &Did{
		masterBytes:  masterBytes,
		agentDid:     agentDid,
		agentPurpose: "agent",
		twinPurpose:  "twin",
	}

	rtn.agentKeyName = rtn.addHashPrefix(agentKeyName)

	return rtn, nil
}

func (d *Did) addHashPrefix(input string) string {
	if !strings.HasPrefix(input, "#") {
		return "#" + input
	}
	return input
}

func (d *Did) CreateAgentDid(keyname string) ([]byte, error) {
	if !strings.HasPrefix(keyname, "#") {
		keyname = "#" + keyname
	}

	doc, err := d.createDid(d.agentPurpose, keyname)
	if err != nil {
		return nil, err
	}
	return []byte(doc.ID), nil
}

func (d *Did) CreateTwinDid(keyname string) ([]byte, error) {
	log.Info("CreateTwinDid")
	purpose := "twin"

	if !strings.HasPrefix(keyname, "#") {
		keyname = "#" + keyname
	}

	twinDoc, err := d.createDid(purpose, keyname)
	if err != nil {
		return nil, err
	}

	_, privateECDSA := d.getPublicPrivate("twin", keyname)

	number := 0
	agentPublicECDSA, agentPrivateECDSA := GetPublicPrivate(d.masterBytes, "agent", uint64(number))
	//agentPublicECDSA, agentPrivateECDSA := d.getPublicPrivate("agent", d.agentKeyName)
	agentID := d.getIdentifier(agentPublicECDSA)

	proof, err := identity.NewProof(twinDoc.ID, agentPrivateECDSA)
	if err != nil {
		log.Printf("Failed to make proof signature")
		return nil, err
	}

	deleg := identity.Delegation{
		ID:         d.agentKeyName,
		Controller: agentID + d.agentKeyName,
		Proof:      proof,
	}

	log.Info(deleg)

	found := false
	for _, dd := range twinDoc.DelegateControl {
		if dd.ID == d.agentKeyName {
			found = true
			break
		}
	}
	if !found {
		//delegDoc.DelegateAuthentication = append(delegDoc.DelegateAuthentication, deleg)
		twinDoc.DelegateControl = append(twinDoc.DelegateControl, deleg)
	}

	// Register the delegation
	audience, err := d.getResolverAudience()
	docClaims := &identity.DIDDocumentClaims{
		Issuer:       twinDoc.ID + twinDoc.PublicKeys[0].ID,
		Audience:     audience,
		Doc:          twinDoc,
		PrivKeyECDSA: privateECDSA,
	}
	if err != nil {
		return nil, err
	}

	err = d.registerDoc(docClaims)
	if err != nil {
		return nil, err
	}

	return []byte(twinDoc.ID), nil
}

// createDid: Attempt to create a DID.  If it exists in the resolver the fetched doc will be returned.
func (d *Did) createDid(purpose string, keyname string) (*identity.DIDDocument, error) {
	log.Infof("createDid [%s] [%s]", purpose, keyname)
	publicECDSA, privateECDSA := d.getPublicPrivate(purpose, keyname)
	id := d.getIdentifier(publicECDSA)

	dtype, _ := identity.StringToDIDType(purpose)
	doc, err := identity.NewDIDDocument(dtype, privateECDSA, keyname)
	if err != nil {
		return nil, err
	}

	audience, err := d.getResolverAudience()
	if err != nil {
		return nil, err
	}
	log.Infof("Audience: %s", audience)
	issID, err := d.joinIdentifierKeyname(id, keyname)
	if err != nil {
		log.Infof("Issuer ID invalid %s %s", issID, err)
		return nil, err
	}
	docClaims := &identity.DIDDocumentClaims{
		Issuer:       issID,
		Audience:     audience,
		Doc:          doc,
		PrivKeyECDSA: privateECDSA,
	}

	rslv, _ := identity.NewResolverClient()
	fetchDoc, err := rslv.Get(id, true)
	if fetchDoc != nil && err == nil {
		log.Infof("Document already exists subject %s", id)
		return fetchDoc, nil
	}

	err = d.registerDoc(docClaims)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// GetIdentifier returns DID Identifier from ecdsa publickey
func (d *Did) getIdentifier(publicECDSA *ecdsa.PublicKey) string {
	pubKeyBytes := identity.ECDSAPublicToBytes(publicECDSA)
	return identity.MakeIdentifier(pubKeyBytes)
}

func (d *Did) getPublicPrivate(purpose string, name string) (*ecdsa.PublicKey, *ecdsa.PrivateKey) {
	privateECDSA, err := identity.NewPrivateKeyECDSAFromPathString(d.masterBytes, purpose, name)
	if err != nil {
		log.Fatal(err)
	}

	publicECDSA, err := identity.ECDSAPrivateToPublic(privateECDSA)
	if err != nil {
		log.Fatal(err)
	}

	return publicECDSA, privateECDSA
}

func (d *Did) joinIdentifierKeyname(id string, name string) (string, error) {
	if err := identity.ValidateIdentifier(id); err != nil {
		return "", err
	}
	if len(name) == 0 {
		return id, nil
	} else {
		if !identity.ValidateName(name) {
			return "", errors.New("Name is not valid")
		}
		result := id + name
		if err := identity.ValidateIdentifier(result); err != nil {
			return "", err
		}
		return result, nil
	}
}

func (d *Did) getResolverAudience() (string, error) {
	addr := ""
	rslv, err := identity.NewResolverClient()
	if err != nil {
		return addr, fmt.Errorf("Failed to setup resolver: %s", err)
	}
	addr, err = rslv.GetAddr()
	if err != nil {
		return addr, fmt.Errorf("Failed to get resolver address: %s", err)
	}
	return addr, nil
}

func (d *Did) registerDoc(docClaims *identity.DIDDocumentClaims) error {
	rslv, err := identity.NewResolverClient()
	if err != nil {
		log.Fatalf("Failed to setup resolver: %s", err)
		return err
	}

	docClaims.Doc.UpdateTime = time.Now().UnixNano() / 1e6 // Unix time millis!

	tkn, err := identity.NewDocumentToken(docClaims)
	if err != nil {
		log.Printf("Failed to create document token: %s", err)
		return err
	}

	fetch, err := rslv.Get(docClaims.Issuer, true)
	if fetch != nil && err == nil {
		log.Printf("Overwriting already existing subject %s", docClaims.Issuer)
	}

	_, err = identity.VerifyDocument(tkn, true)
	log.Printf("registerDoc identity.VerifyDocument %v", err)

	return rslv.Register(tkn)
}

const maxDuration = 60 * 60 * 24

func (d *Did) GetAgentJwt(userDid string, duration string) (string, error) {

	purpose := "agent"
	audience := "not-currently-used.example.com"
	issuerID := fmt.Sprintf("%s%s", d.agentDid, d.agentKeyName)

	flagAudience := audience
	flagSubject := userDid
	flagIssuer := d.agentDid
	flagKeyname := d.agentKeyName
	flagDurationString := duration
	flagDuration, err := time.ParseDuration(flagDurationString)
	if err != nil {
		log.Error(err)
		return "", err
	}

	number := 0
	publicECDSA, privateECDSA := GetPublicPrivate(d.masterBytes, purpose, uint64(number))

	if purpose != "agent" && purpose != "twin" {
		return "", fmt.Errorf("purpose must be 'agent' or 'twin'")
	}

	if flagDuration.Seconds() <= 0 || flagDuration.Seconds() > maxDuration {
		return "", fmt.Errorf("duration must be >0 and <%d", maxDuration)
	}

	if len(flagAudience) == 0 {
		return "", fmt.Errorf("audience is required")
	}
	if identity.ValidateIdentifier(flagSubject) != nil {
		return "", fmt.Errorf("subject is required")
	}

	rslv, err := identity.NewResolverClient()
	if err != nil {
		return "", fmt.Errorf("Failed to setup resolver: %s", err)
	}
	flagIssuerDoc, err := rslv.Get(flagIssuer, true)
	if err != nil {
		return "", fmt.Errorf("Failed to get Issuer %s Document %s", flagIssuer, err)
	}

	// Check seed+purpose allowed to work on Issuers behalf
	var issuerKey *identity.IssuerKey

	issuerID, err = d.joinIdentifierKeyname(flagIssuer, flagKeyname)
	if err != nil {
		return "", fmt.Errorf("issuer invalid %s %s", issuerID, err)
	}

	issuerKey, err = identity.FindIssuerByDocument(flagIssuerDoc, issuerID, false, true)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	if !issuerKey.Matches(publicECDSA) {
		return "", fmt.Errorf("Issuer %s does not match key (seed + purpose + number)", issuerID)
	}

	// Check issuer allowed to auth as Subject!
	flagSubjectDoc, err := rslv.Get(flagSubject, true)
	if err != nil {
		return "", fmt.Errorf("Failed to get Subject %s Document %s", flagSubject, err)
	}
	_, err = identity.FindIssuerByDocument(flagSubjectDoc, issuerID, false, true)
	if err != nil {
		return "", fmt.Errorf("%s not allowed to authenticate as %s", issuerID, flagSubject)
	}

	authReq := &identity.AuthenticationClaims{
		Issuer:       issuerKey.IssuerID,
		Subject:      flagSubject,
		Audience:     flagAudience,
		Duration:     flagDuration,
		PrivKeyECDSA: privateECDSA,
	}
	tkn, err := identity.NewAuthenticationToken(authReq)
	if err != nil {
		log.Fatalf("Unable to make auth token: %s", err)
	}

	//	log.Printf(tkn)
	return tkn, nil
}

func GetPublicPrivate(masterBytes []byte, purpose string, number uint64) (*ecdsa.PublicKey, *ecdsa.PrivateKey) {
	privateECDSA, err := identity.NewPrivateKeyECDSAFromPath(masterBytes, purpose, number)
	if err != nil {
		log.Fatal(err)
	}

	publicECDSA, err := identity.ECDSAPrivateToPublic(privateECDSA)
	if err != nil {
		log.Fatal(err)
	}

	return publicECDSA, privateECDSA
}
