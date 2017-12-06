package timesheet

const Template = `{{range .}}{{.Start}}:{{.End}} {{print .Msg}}
{{end}}`
