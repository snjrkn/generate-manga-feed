package generatefeed

import "github.com/snjrkn/generate-manga-feed/internal/service"

func KurageFarm() (string, error) {
	return service.KurageFarm().MakeFeed()
}

func ComicdaysOneshot() (string, error) {
	return service.ComicDaysOneshot().MakeFeed()
}

func ComicdaysNewcomer() (string, error) {
	return service.ComicDaysNewcomer().MakeFeed()
}

func AndSofa() (string, error) {
	return service.AndSofa().MakeFeed()
}

func Toti() (string, error) {
	return service.Toti().MakeFeed()
}

func Matogrosso() (string, error) {
	return service.Matogrosso().MakeFeed()
}

func KurageAward() (string, error) {
	return service.KurageAward().MakeFeed()
}

func ComicEssayGekijo() (string, error) {
	return service.ComicEssayGekijo().MakeFeed()
}

func ComiplexOneshot() (string, error) {
	return service.ComiplexOneshot().MakeFeed()
}

func ComicBoostOneshot() (string, error) {
	return service.ComicBoostOneshot().MakeFeed()
}

func YoungAnimalOneshot() (string, error) {
	return service.YoungAnimalOneshot().MakeFeed()
}

func ComicEssayContest() (string, error) {
	return service.ComicEssayContest().MakeFeed()
}

func ComicBunchKaiAward() (string, error) {
	return service.ComicBunchKaiAward().MakeFeed()
}

func ComicBunchKaiOneshot() (string, error) {
	return service.ComicBunchKaiOneshot().MakeFeed()
}

func AfternoonAward() (string, error) {
	return service.AfternoonAward().MakeFeed()
}

func ShonenMagazineAward() (string, error) {
	return service.ShonenMagazineAward().MakeFeed()
}

func ShonenMagazineRise() (string, error) {
	return service.ShonenMagazineRise().MakeFeed()
}
