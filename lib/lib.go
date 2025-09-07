package lib

import "github.com/snjrkn/generate-manga-feed/internal/service"

func KurageFarm() (string, error) {
	return service.KurageFarm().MakeFeed()
}

func ComicdaysOneshot() (string, error) {
	return service.ComicdaysOneshot().MakeFeed()
}

func ComicdaysNewcomer() (string, error) {
	return service.ComicdaysNewcomer().MakeFeed()
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

func ComicEssayAward() (string, error) {
	return service.ComicEssayAward().MakeFeed()
}
