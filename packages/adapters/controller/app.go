package controller

type AppController struct {
	Search    interface{ Search }
	User      interface{ User }
	Question  interface{ Question }
	Selection interface{ Selection }
}
