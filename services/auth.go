package services

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var SERVICE_URL = os.Getenv("AUTH_SERVICE_URL")

type AuthService interface {
	Login(service string, email string, password string) (string, error)
	Check(service string, token string) (string, error)
	Update() (bool, error)
	Create(service string, email string, password string) (bool, error)
	Delete(string) (bool, error)
}

type userAuthService struct {
	token string
}

func NewUserAuthService() AuthService {
	return &userAuthService{
		token: "",
	}
}

func (userAuth *userAuthService) Login(service string, email string, password string) (string, error) {
	resp, err := http.PostForm(SERVICE_URL+"/login",
		url.Values{"service": {service}, "email": {email}, "password": {password}})
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	} else {
		return "", err
	}

}

func (userAuth *userAuthService) Check(service string, token string) (string, error) {
	resp, err := http.Get(SERVICE_URL + "/user?service=" + service + "&token=" + token)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	} else {
		return "", err
	}
}

func (userAuth *userAuthService) Update() (bool, error) {
	return false, nil
}

func (userAuth *userAuthService) Create(service string, email string, password string) (bool, error) {
	resp, err := http.PostForm(SERVICE_URL+"/user",
		url.Values{"service": {service}, "email": {email}, "password": {password}})
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}

}

func (userAuth *userAuthService) Delete(string) (bool, error) {
	return false, nil
}
