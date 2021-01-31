# IndexOfBrowser Library

Takes those "index of" directories and provides methods to explore them.

```
func main() {
	listFilesRecurse(
		IndexOfBrowser.NewBrowser(os.Args[1]))
	//find all trs, ingore
}

func listFilesRecurse(browser *IndexOfBrowser.Browser) {
	files, err := browser.Ls()
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range files  {
		if entry.Directory {
			browser.Pushd(browser.Cwd() + entry.Name)
			listFilesRecurse(browser)
			browser.Popd()
		} else {
			log.Printf("File %s", entry.Url)
		}
	}
}
```