package generatemangafeed

import (
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/service"
)

func KurageFarm(string) (string, error) {
	return generator.NewGenerator(service.AfternoonAward()).MakeFeed()
}

func ComicdaysOneshot(string) (string, error) {
	return generator.NewGenerator(service.ComicDaysOneshot()).MakeFeed()
}

func ComicdaysNewcomer(string) (string, error) {
	return generator.NewGenerator(service.ComicDaysNewcomer()).MakeFeed()
}

func AndSofa(string) (string, error) {
	return generator.NewGenerator(service.AndSofa()).MakeFeed()
}

func Toti(string) (string, error) {
	return generator.NewGenerator(service.Toti()).MakeFeed()
}

func Matogrosso(string) (string, error) {
	return generator.NewGenerator(service.Matogrosso()).MakeFeed()
}

func KurageBunchAward(string) (string, error) {
	return generator.NewGenerator(service.KurageBunchAward()).MakeFeed()
}

func ComicEssayGekijo(string) (string, error) {
	return generator.NewGenerator(service.ComicEssayGekijo()).MakeFeed()
}

func ComiplexOneshot(string) (string, error) {
	return generator.NewGenerator(service.ComiplexOneshot()).MakeFeed()
}

func ComicBoostOneshot(string) (string, error) {
	return generator.NewGenerator(service.ComicBoostOneshot()).MakeFeed()
}

func YoungAnimalOneshot(string) (string, error) {
	return generator.NewGenerator(service.YoungAnimalOneshot()).MakeFeed()
}

func ComicEssayContest(string) (string, error) {
	return generator.NewGenerator(service.ComicEssayContest()).MakeFeed()
}

func ComicBunchKaiAward(string) (string, error) {
	return generator.NewGenerator(service.ComicBunchKaiAward()).MakeFeed()
}

func ComicBunchKaiOneshot(string) (string, error) {
	return generator.NewGenerator(service.ComicBunchKaiOneshot()).MakeFeed()
}

func AfternoonAward(string) (string, error) {
	return generator.NewGenerator(service.AfternoonAward()).MakeFeed()
}

func ShonenMagazineAward(string) (string, error) {
	return generator.NewGenerator(service.ShonenMagazineAward()).MakeFeed()
}

func ShonenMagazineRise(string) (string, error) {
	return generator.NewGenerator(service.ShonenMagazineRise()).MakeFeed()
}

func ChampionCrossOneshot(string) (string, error) {
	return generator.NewGenerator(service.ChampionCrossOneshot()).MakeFeed()
}

func KurageBunchOneshot(string) (string, error) {
	return generator.NewGenerator(service.KurageBunchOneshot()).MakeFeed()
}

func ComicActionOneshot(string) (string, error) {
	return generator.NewGenerator(service.ComicActionOneshot()).MakeFeed()
}

func ComicBoostRensai(productId string) (string, error) {
	return generator.NewGenerator(service.ComicBoostRensai(productId)).MakeFeed()
}
