package cron

//支持秒级任务

func NewWithSeconds() *Cron {

	secondParser := NewParser(Second | Minute | Hour |
		Dom | Month | DowOptional | Descriptor)

	return New(WithParser(secondParser), WithChain())

}
