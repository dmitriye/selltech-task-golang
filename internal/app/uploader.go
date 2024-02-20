package app

import (
	"encoding/xml"
	"fmt"
	"sync"
)

type Uploader struct {
	wg         sync.WaitGroup
	appState   *AppState
	repo       *EntryRepository
	ofacApi    *OfacApi
	workersNum int
}

func NewUploader(s *AppState, repo *EntryRepository, o *OfacApi, n int) *Uploader {
	return &Uploader{
		appState:   s,
		repo:       repo,
		ofacApi:    o,
		workersNum: n,
	}
}

type xmlEntry struct {
	UID       int    `xml:"uid"`
	FirstName string `xml:"firstName"`
	LastName  string `xml:"lastName"`
	Type      string `xml:"sdnType"`
}

func (u *Uploader) Run() error {
	const fn = "app.uploader.Run"

	if u.appState.IsUpdating() {
		return fmt.Errorf("incorrect application state: %s", u.appState.GetName())
	}

	u.appState.SetState(S_UPDATING)

	if err := u.internalRun(); err != nil {
		u.appState.SetState(S_EMPTY)
		u.repo.ClearAll()
		return fmt.Errorf("%s :%w", fn, err)
	}

	u.appState.SetState(S_OK)

	return nil
}

func (u *Uploader) startWorkers(n int, ch <-chan xmlEntry, h func(e xmlEntry)) {
	for i := 0; i < n; i++ {
		u.wg.Add(1)
		go func() {
			defer u.wg.Done()
			for v := range ch {
				h(v)
			}
		}()
	}
}

func (u *Uploader) internalRun() error {
	const fn = "app.uploader.internalRun"

	ch := make(chan xmlEntry, u.workersNum)

	handle, err := u.ofacApi.GetSdnEntries()
	if err != nil {
		return fmt.Errorf("%s :%w", fn, err)
	}
	defer handle.Close()

	decoder := xml.NewDecoder(handle)

	err = u.repo.ClearAll()
	if err != nil {
		return fmt.Errorf("%s :%w", fn, err)
	}

	u.startWorkers(u.workersNum, ch, func(e xmlEntry) {
		u.repo.CreateOrUpdate(Entry{
			UID:       e.UID,
			FirstName: e.FirstName,
			LastName:  e.LastName,
		})
	})

	var entry xmlEntry
	for {
		t, _ := decoder.Token()

		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "sdnEntry" {
				decoder.DecodeElement(&entry, &se)

				if entry.Type == "Individual" {
					// fmt.Printf("%#v\n", entry)
					ch <- entry
				}
			}
		}
	}

	close(ch)
	u.wg.Wait()

	return nil
}
