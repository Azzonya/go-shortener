package multicheker

import (
	"reflect"
	"testing"
)

func Test_getStaticCheckAnalyzers(t *testing.T) {
	checkIDs := []string{
		"SA4013", "SA4015", "SA1023", "SA1026", "SA1030", "SA3001", "SA4021", "SA9008",
		"SA1002", "SA1007", "SA5012", "SA6003", "SA9001", "SA9006", "SA1012", "SA1013",
		"SA3000", "SA4010", "SA4016", "SA4028", "SA5002", "SA5007", "SA1001", "SA1008",
		"SA9007", "SA6001", "SA9003", "SA4029", "SA5009", "SA5011", "SA1021", "SA1027",
		"SA5004", "SA4011", "SA4012", "SA4001", "SA4006", "SA4014", "SA4026", "SA5008",
		"SA1016", "SA1029", "SA1024", "SA1025", "SA4005", "SA4019", "SA1004", "SA1006",
		"SA4017", "SA4030", "SA5000", "SA5010", "SA6005", "SA9004", "SA1011", "SA1015",
		"SA5003", "SA1014", "SA1019", "SA4003", "SA4018", "SA9002", "SA9005", "SA1020",
		"SA2001", "SA5005", "SA4000", "SA4027", "SA1018", "SA1028", "SA4008", "SA4009",
		"SA4020", "SA4024", "SA1000", "SA1010", "SA4025", "SA2002", "SA6000", "SA6002",
		"SA1003", "SA1017", "SA4023", "SA5001", "SA2003", "SA4022", "SA4004", "SA4031",
		"SA1005", "SA2000", "S1006", "S1008", "S1025", "S1024", "S1029", "S1034", "S1002",
		"S1012", "S1019", "S1020", "S1033", "S1018", "S1021", "S1031", "S1040", "S1001",
		"S1009", "S1016", "S1017", "S1000", "S1005", "S1011", "S1030", "S1003", "S1004",
		"S1028", "S1039", "S1007", "S1010", "S1023", "S1038", "S1032", "S1035", "S1036",
		"S1037", "QF1002", "QF1006", "QF1012", "QF1010", "QF1001", "QF1003", "QF1004",
		"QF1005", "QF1007", "QF1008", "QF1009", "QF1011", "ST1005", "ST1011", "ST1013",
		"ST1017", "ST1019", "ST1022", "ST1000", "ST1006", "ST1012", "ST1015", "ST1020",
		"ST1001", "ST1003", "ST1008", "ST1016", "ST1021", "ST1018", "ST1023",
	}

	tests := []struct {
		name string
		want int
	}{
		{
			name: "static check analyzers",
			want: len(checkIDs),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStaticCheckAnalyzers(); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getStaticCheckAnalyzers() = %v, want %v", got, tt.want)
			}
		})
	}
}
