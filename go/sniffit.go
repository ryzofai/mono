func main() {

	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		fmt.Println(err)
	}
  fmt.Println(configuration.File_pattern)
	re := regexp.MustCompile(configuration.File
  fmt.Println(re.MatchString("pas5-h2h-2018-05-15-111037.err"))
}
