package generatemangafeed

import "github.com/snjrkn/generate-manga-feed/internal/service"

func KurageFarm(string) (string, error) {
	return service.KurageFarm().MakeFeed()
}

func ComicdaysOneshot(string) (string, error) {
	return service.ComicDaysOneshot().MakeFeed()
}

func ComicdaysNewcomer(string) (string, error) {
	return service.ComicDaysNewcomer().MakeFeed()
}

func AndSofa(string) (string, error) {
	return service.AndSofa().MakeFeed()
}

func Toti(string) (string, error) {
	return service.Toti().MakeFeed()
}

func Matogrosso(string) (string, error) {
	return service.Matogrosso().MakeFeed()
}

func KurageBunchAward(string) (string, error) {
	return service.KurageBunchAward().MakeFeed()
}

func ComicEssayGekijo(string) (string, error) {
	return service.ComicEssayGekijo().MakeFeed()
}

func ComiplexOneshot(string) (string, error) {
	return service.ComiplexOneshot().MakeFeed()
}

func ComicBoostOneshot(string) (string, error) {
	return service.ComicBoostOneshot().MakeFeed()
}

func YoungAnimalOneshot(string) (string, error) {
	return service.YoungAnimalOneshot().MakeFeed()
}

func ComicEssayContest(string) (string, error) {
	return service.ComicEssayContest().MakeFeed()
}

func ComicBunchKaiAward(string) (string, error) {
	return service.ComicBunchKaiAward().MakeFeed()
}

func ComicBunchKaiOneshot(string) (string, error) {
	return service.ComicBunchKaiOneshot().MakeFeed()
}

func AfternoonAward(string) (string, error) {
	return service.AfternoonAward().MakeFeed()
}

func ShonenMagazineAward(string) (string, error) {
	return service.ShonenMagazineAward().MakeFeed()
}

func ShonenMagazineRise(string) (string, error) {
	return service.ShonenMagazineRise().MakeFeed()
}

func ChampionCrossOneshot(string) (string, error) {
	return service.ChampionCrossOneshot().MakeFeed()
}

func KurageBunchOneshot(string) (string, error) {
	return service.KurageBunchOneshot().MakeFeed()
}

func ComicActionOneshot(string) (string, error) {
	return service.ComicActionOneshot().MakeFeed()
}

func ComicBoostRensai(productId string) (string, error) {
	return service.ComicBoostRensai(productId).MakeFeed()
}
