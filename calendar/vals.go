package calendar

var days = [7]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var months = [12]string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

var monthNdays = map[string][]int{
	months[0]:  {31},
	months[1]:  {28, 29},
	months[2]:  {31},
	months[3]:  {30},
	months[4]:  {31},
	months[5]:  {30},
	months[6]:  {31},
	months[7]:  {31},
	months[8]:  {30},
	months[9]:  {31},
	months[10]: {30},
	months[11]: {31},
}
