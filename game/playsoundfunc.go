package game

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Pour jouer les différents sons et les chargers
const sampleRate = 44100

var sound []sound_object

type sound_object struct {
	data    *audio.Player // son
	name    string        // nom son
	stream  io.Closer     // lecteur
	started bool          // il y a un temps d'attente avant que le son soit considérer comme 'joué actuellement' donc tant que started n'est pas mis à vrai, alors c'est que le son charge
}

var aud *audio.Context = audio.NewContext(sampleRate)

// Pour le sound manager
var ambiance_ordre_de_passage []int = []int{1, 2, 3, 4, 5, 6, 7, 8}
var ambiance_lecteur_position int = rand.Intn(7) + 1

func playsound(fileRep, name string) {
	// Supprimez cette ligne (redéclaration inutile)
	// var aud *audio.Context
	// aud = audio.NewContext(sampleRate)

	file, err := os.Open(fileRep)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture du fichier: %v", err)
	}

	streamer, err := wav.Decode(aud, file) // Utilisez le contexte global `aud`
	if err != nil {
		log.Fatalf("Erreur lors du décodage du fichier audio: %v", err)
	}

	player, err := aud.NewPlayer(streamer)
	if err != nil {
		log.Fatalf("Erreur lors de la création du player audio: %v", err)
	}

	sound = append(sound, sound_object{
		data:   player,
		name:   name,
		stream: file,
	})
	fmt.Printf("Son '%s' ajouté et joué.\n", name)

	player.Rewind()
	player.Play()
}

func stopsound(name string) {
	for i := len(sound) - 1; i >= 0; i-- {
		if sound[i].name == name {
			// Fermez le player
			sound[i].data.Close()

			// Fermez explicitement le flux si disponible
			if sound[i].stream != nil {
				sound[i].stream.Close()
			}

			// Supprimez l'objet de la liste
			sound = append(sound[:i], sound[i+1:]...)
			fmt.Printf("Son '%s' arrêté et supprimé.\n", name)
			return
		}
	}
	fmt.Printf("Son '%s' introuvable.\n", name)
}

func isplayingsound(name string) bool {
	// Cette fonction permet de savoir si le son << name >> est joué actuellement

	z := len(sound)
	for index := 0; index < len(sound); index++ {
		if !sound[index].data.IsPlaying() {
			if sound[index].started {
				sound[index].data.Close()
				if sound[index].stream != nil {
					sound[index].stream.Close()
				}
				sound = append(sound[:index], sound[index+1:]...)
				z = z - 1
			}
		} else {
			if sound[index].name == name {
				return true
			}
			sound[index].started = true
		}
	}
	return false
}
func changevolumesound(name string, volume float64) {
	// Cette fonction permet de changer le volume du/des son nommé << name >>

	z := len(sound)
	for index := 0; index < z; index++ {
		if sound[index].name == name {
			sound[index].data.SetVolume(volume)
		}
	}
}

func (g *Game) manage_sounds(moving bool) {
	// S'occupe des musiques d'ambiance, gestion intelligente des sons de façon à ne pas ennuyer le joueur avec une musique jouée trop de fois

	// jouer les musiques
	if !isplayingsound("ambiance") {
		playsound("../assets/sounds/musics/C418_"+strconv.Itoa(ambiance_ordre_de_passage[ambiance_lecteur_position])+".wav", "ambiance")
		if rand.Intn(2) == 0 {
			ambiance_ordre_de_passage = append(ambiance_ordre_de_passage[1:], ambiance_ordre_de_passage[0])
		} else {
			br := ambiance_ordre_de_passage[7]
			ambiance_ordre_de_passage = append(append(ambiance_ordre_de_passage[1:7], ambiance_ordre_de_passage[0]), br)
		}
	}

	// Petit ajout qui baisse un peu le son des musiques d'ambiance dans les ui (sauf le menu)
	if g.PlayerOnUi == "" || g.PlayerOnUi == "menu" {
		changevolumesound("ambiance", 1)
	} else {
		changevolumesound("ambiance", 0.5)
	}

	// S'occupe des bruits de pas
	if moving {
		g.manage_footsteps()
	}

	// S'occuper des sons que peuvent produire les blocs
	g.manage_sound_from_blocs()
}

func list_contain_x(list []int, x int) bool {
	// Vérifie si la liste <list> contient bien le nombre <x>
	z := len(list)
	for i := 0; i < z; i++ {
		if list[i] == x {
			return true
		}
	}
	return false
}

func (g *Game) manage_footsteps() {
	// Cette fonction s'occupe des bruits de pas, selon sur quel type de bloc on marche, alors un des 4 bruits du bloc sera joué

	x, y := g.Character.GetContentCharacterPos(g.Camera.X, g.Camera.Y)
	blocID := g.Floor.GetContent()[y][x]
	sound_type := "grass"
	if list_contain_x(configuration.Global.WoodBlocks, blocID) {
		sound_type = "wood"
	} else if list_contain_x(configuration.Global.StoneBlocks, blocID) {
		sound_type = "stone"
	} else if list_contain_x(configuration.Global.SandBlocks, blocID) {
		sound_type = "sand"
	}

	if !isplayingsound("Footsteps") {
		playsound("../assets/sounds/effects/walk/"+sound_type+"/#"+strconv.Itoa(rand.Intn(3)+1)+".wav", "Footsteps")
	}
}

func (g *Game) getClosestDistancePlayerToBlocIDinContent(blocID int) int {
	// Cette fonction permet de connaître la distance sur l'axe X entre un objet et le joueur
	//(il prend l'objet le plus proche)

	content_X, content_Y := g.Character.GetContentCharacterPos(g.Camera.X, g.Camera.Y)
	content := g.Floor.GetContent()
	closest_distance_from_player := -1
	var br_cd int
	for Y := 0; Y < configuration.Global.NumTileY; Y++ {
		for X := 0; X < configuration.Global.NumTileX; X++ {
			if content[Y][X] == blocID {
				if content_X > X {
					br_cd = content_X - X
				} else {
					br_cd = X - content_X
				}
				if content_Y > Y {
					br_cd += content_Y - Y
				} else {
					br_cd += Y - content_Y
				}
				if closest_distance_from_player == -1 || br_cd < closest_distance_from_player {
					closest_distance_from_player = br_cd
				}
			}
		}
	}
	return closest_distance_from_player
}

func (g *Game) manage_sound_from_blocs() {
	// Cette fonction joue les bruits produits par les blocs

	// Bruit pour les fours
	distance := g.getClosestDistancePlayerToBlocIDinContent(47)
	if distance > -1 {
		if !isplayingsound("Furnace") {
			playsound("../assets/sounds/effects/blocs/furnace.wav", "Furnace")
		}
		if distance > 7 {
			distance = 7
		}
		changevolumesound("Furnace", 1.0-float64(distance)*(1.0/7.0))
	} else {
		changevolumesound("Furnace", 0.0)
	}
}
