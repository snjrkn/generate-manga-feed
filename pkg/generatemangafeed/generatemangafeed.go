package generatemangafeed

import (
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/service"
)

func KurageFarm(string) (string, error) {
	return generator.NewGenerator(service.NewKurageFarm()).MakeFeed()
}

func ComicdaysOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewComicDaysOneshot()).MakeFeed()
}

func ComicdaysNewcomer(string) (string, error) {
	return generator.NewGenerator(service.NewComicDaysNewcomer()).MakeFeed()
}

func AndSofa(string) (string, error) {
	return generator.NewGenerator(service.NewAndSofa()).MakeFeed()
}

func Toti(productId string) (string, error) {
	return generator.NewGenerator(service.NewToti(productId)).MakeFeed()
}

func Matogrosso(string) (string, error) {
	return generator.NewGenerator(service.NewMatogrosso()).MakeFeed()
}

func KurageBunchAward(string) (string, error) {
	return generator.NewGenerator(service.NewKurageBunchAward()).MakeFeed()
}

func ComicEssayGekijo(string) (string, error) {
	return generator.NewGenerator(service.NewComicEssayGekijo()).MakeFeed()
}

func ComiplexOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewComiplexOneshot()).MakeFeed()
}

func ComicBoostOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewComicBoostOneshot()).MakeFeed()
}

func YoungAnimalOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewYoungAnimalOneshot()).MakeFeed()
}

func ComicEssayContest(string) (string, error) {
	return generator.NewGenerator(service.NewComicEssayContest()).MakeFeed()
}

func ComicBunchKaiAward(string) (string, error) {
	return generator.NewGenerator(service.NewComicBunchKaiAward()).MakeFeed()
}

func ComicBunchKaiOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewComicBunchKaiOneshot()).MakeFeed()
}

func AfternoonAward(string) (string, error) {
	return generator.NewGenerator(service.NewAfternoonAward()).MakeFeed()
}

func ShonenMagazineAward(string) (string, error) {
	return generator.NewGenerator(service.NewShonenMagazineAward()).MakeFeed()
}

func ShonenMagazineRise(string) (string, error) {
	return generator.NewGenerator(service.NewShonenMagazineRise()).MakeFeed()
}

func ChampionCrossOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewChampionCrossOneshot()).MakeFeed()
}

func KurageBunchOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewKurageBunchOneshot()).MakeFeed()
}

func ComicActionOneshot(string) (string, error) {
	return generator.NewGenerator(service.NewComicActionOneshot()).MakeFeed()
}

func ComicBoostRensai(contentId string) (string, error) {
	return generator.NewGenerator(service.NewComicBoostRensai(contentId)).MakeFeed()
}
