package main

var gIService IService

type IService interface {
	Start() (err error)
	Stop() (err error)
}
