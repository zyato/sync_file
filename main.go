package main

func main() {
	config := getConfig()
	logger := getLogger()
	err := newSyncHelper(config.SourceDir, config.TargetDir).
		withInterval(config.Interval).
		withLogger(logger).
		run()
	panic("初始化syncHelper失败：" + err.Error())
}
