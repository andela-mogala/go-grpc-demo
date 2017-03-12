package main

import "github.com/andela-mogala/go-grpc-demo/src/grpcdemo/pb"

var Employees = []pb.Employee{
	pb.Employee{
		Id:                  1,
		BadgeNumber:         2080,
		FirstName:           "Grace",
		LastName:            "Decker",
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
	pb.Employee{
		Id:                  2,
		BadgeNumber:         7538,
		FirstName:           "Amity",
		LastName:            "Fuller",
		VacationAccrualRate: 2.3,
		VacationAccrued:     23.4,
	},
	pb.Employee{
		Id:                  3,
		BadgeNumber:         5144,
		FirstName:           "Keaton",
		LastName:            "Willis",
		VacationAccrualRate: 3,
		VacationAccrued:     31.7,
	},
}
