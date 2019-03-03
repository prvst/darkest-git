package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

var baseURl = "https://gamepedia.cursecdn.com/darkestdungeon_gamepedia/"

// Config
var quoteConfigAudio = [...]string{
	"e/e2/Vo_narr_tut_firstprov.ogg",
}
var quoteConfigText = [...]string{
	"The cost of preparedness - measured now in gold, later in blood",
}

// Init
var quoteInitAudio = [...]string{
	"3/30/Vo_narr_tut_firsttown.ogg",
	"b/b5/Vo_narr_tut_firstquest.ogg",
	"8/80/Vo_narr_tut_firststage.ogg",
	"2/29/Vo_narr_town_dismiss_06.ogg",
}
var quoteInitText = [...]string{
	"Welcome home, such as it is. This squalid hamlet, these corrupted lands, they are yours now, and you are bound to them",
	"This sprawling estate, a Mecca of madness and morbidity. Your work begins...",
	"Women and men; soldiers and outlaws; fools and corpses. All will find their way to us now that the road is clear",
	"The task ahead is terrible, and weakness cannot be tolerated.",
}

// rm
var quoteRmAudio = [...]string{
	"d/d0/Vo_narr_tut_firstgrave.ogg",
	"2/29/Vo_narr_town_upgrade_guild_01.ogg",
	"a/a1/Vo_narr_town_dismiss_05.ogg",
	"e/ec/Vo_narr_town_dismiss_08.ogg",
}
var quoteRmText = [...]string{
	"Most will end up here, covered in the poisoned earth, awaiting merciful oblivion",
	"Some may fall, but their knowledge lives on",
	"This one has become vestigial, useless",
	"Slumped shoulders, wild eyes, and a stumbling gait - this one is no more good to us",
}

// log
var quoteLogAudio = [...]string{
	"e/e8/Vo_narr_town_backstory_14.ogg",
	"8/80/Vo_narr_town_backstory_02.ogg",
	"9/9f/Vo_narr_neut_weald_gather.ogg",
}
var quoteLogText = [...]string{
	"In time, you will know the tragic extent of my failings...",
	"I remember days when the sun shone, and laughter could be heard from the tavern",
	"Our land is remote and unneighbored. Every lost resource must be recovered.",
}

// reset
var quoteResetAudio = [...]string{
	"4/41/Vo_narr_town_backstory_03.ogg",
	"5/5d/Vo_narr_town_backstory_04.ogg",
	"0/02/Vo_narr_town_backstory_22.ogg",
}
var quoteResetText = [...]string{
	"I was lord of this place, before the crows and rats made it their domain",
	"I see something long-absent in the sunken faces of passersby - a glimmer of hope",
	"An eternity of futile struggle — a penance for my unspeakable transgressions",
}

// show
var quoteShowAudio = [...]string{
	"5/52/Vo_narr_town_backstory_10.ogg",
}
var quoteShowText = [...]string{
	"I can still see their angry faces as they stormed the manor, but I was dead before they found me, and the letter was on its way",
}

// merge
var quoteMergeAudio = [...]string{
	"3/35/Vo_narr_town_backstory_17.ogg",
}
var quoteMergeText = [...]string{
	"We dug for months, years — an eternity. And we were rewarded with madness",
}

// checkout
var quoteCheckoutAudio = [...]string{
	"3/32/Vo_narr_town_backstory_20.ogg",
	"3/38/Vo_narr_town_backstory_19.ogg",
}
var quoteCheckoutText = [...]string{
	"Can you feel it? The walls between the sane world and that unplumbed dimension of delirium are tenuously thin here...",
	"All the decadent horrors I have seen pale in comparison with that final, crowning thing. I could not look, nor could I look away!",
}

// stash
var quoteStashAudio = [...]string{
	"b/b0/Vo_narr_town_upgrade_blacksmith_02.ogg",
}
var quoteStashText = [...]string{
	"Fan the flames! Mold the metal! We are raising an army!",
}

// add
var quoteAddAudio = [...]string{
	"8/8a/Vo_narr_tut_firstnomad.ogg",
	"2/21/Vo_narr_town_upgrade_stage_03.ogg",
	"7/72/Vo_narr_neut_crypts_gather.ogg",
	"6/62/Vo_narr_good_crit_08.ogg",
	"a/a2/Vo_narr_good_crit_09.ogg",
}
var quoteAddText = [...]string{
	"Trinkets and charms, gathered from all the forgotten corners of the earth...",
	"More arrive, foolishly seeking fortune and glory in this domain of the damned",
	"There is power in symbols. Collect the scattered scraps of faith and give comfort to the masses.",
	"Impressive!",
	"The ground quakes!",
}

// commit
var quoteCommitAudio = [...]string{
	"9/9b/Vo_narr_tut_firstbattle.ogg",
	"e/ed/Vo_narr_town_upgrade_tavern_02.ogg",
	"3/3c/Vo_narr_good_healcrit_01.ogg.",
	"b/b2/Vo_narr_good_crit_07.ogg",
}
var quoteCommitText = [...]string{
	"An ambush! Send these vermin a message: the rightful owner has returned, and their kind is no longer welcome",
	"Strong drink, a game of chance, and companionship. The rush of life",
	"Surgical precision!",
	"Precision and power!",
}

// push
var quotePushAudio = [...]string{
	"b/b8/Vo_narr_town_dismiss_07.ogg",
	"2/26/Vo_narr_good_cleanse_crypts_02.ogg",
	"4/41/Narration_activateruins1.ogg",
	"3/32/Vo_narr_good_kill_weak_04.ogg",
	"d/df/Vo_narr_good_crit_06.ogg",
}
var quotePushText = [...]string{
	"It is done. Turn yourself now to the conditions of those poor devils who remain.",
	"Room by room, hall by hall, we reclaim what is ours.",
	"This day belongs to the Light!",
	"Back to the pit!",
	"Masterfully executed!",
}

// pull
var quotePullAudio = [...]string{
	"7/7b/Vo_narr_good_victoryfirst_08.ogg",
	"7/70/Vo_narr_good_victoryfirst_04.ogg",
	"0/04/Vo_narr_good_victoryfirst_09.ogg",
	"2/23/Vo_narr_good_victoryfirst_02.ogg",
	"a/a3/Vo_narr_good_victoryfirst_06.ogg",
}
var quotePullText = [...]string{
	"Seize this momentum! Push on to the task's end!",
	"This expedition, at least, promises success.",
	"Success so clearly in view... or is it merely a trick of the light?",
	"Remind yourself that overconfidence is a slow and insidious killer.",
	"Be wary - triumphant pride precipitates a dizzying fall...",
}

func playAudio(ogg string) {

	// Open first sample File
	f, err := os.Open(ogg)
	if err != nil {
		log.Fatal(err)
	}

	// Decode the .ogg File, if you have a .wav file, use wav.Decode(f)
	s, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Channel, which will signal the end of the playback.
	playing := make(chan struct{})

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		// Callback after the stream Ends
		close(playing)
	})))

	<-playing
}

func downloadFromUrl(url, name string) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(name)
	if err != nil {
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return

	return
}

func runGit(args []string) error {

	cmd := exec.Command("git")

	for _, i := range args {
		if i != "git" {
			cmd.Args = append(cmd.Args, i)
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		return e
	}
	_ = cmd.Wait()

	return nil
}

func main() {

	var args []string
	args = append(args, "git")
	args = append(args, os.Args[1:]...)

	var url string
	var text string
	var name []string
	var index int

	rand.Seed(time.Now().UnixNano())

	// get the file name and build the url path
	if strings.Contains(args[1], "config") {

		index = rand.Intn(len(quoteConfigAudio))
		name = strings.Split(quoteConfigAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteConfigAudio[index])
		text = quoteConfigText[index]

	} else if strings.Contains(args[1], "init") {

		index = rand.Intn(len(quoteInitAudio))
		name = strings.Split(quoteInitAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteInitAudio[index])
		text = quoteInitText[index]

	} else if strings.Contains(args[1], "add") {

		index = rand.Intn(len(quoteAddAudio))
		name = strings.Split(quoteAddAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteAddAudio[index])
		text = quoteAddText[index]

	} else if strings.Contains(args[1], "commit") {

		index = rand.Intn(len(quoteCommitAudio))
		name = strings.Split(quoteCommitAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteCommitAudio[index])
		text = quoteCommitText[index]

	} else if strings.Contains(args[1], "rm") {

		index = rand.Intn(len(quoteRmAudio))
		name = strings.Split(quoteRmAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteRmAudio[index])
		text = quoteRmText[index]

	} else if strings.Contains(args[1], "log") {

		index = rand.Intn(len(quoteLogAudio))
		name = strings.Split(quoteLogAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteLogAudio[index])
		text = quoteLogText[index]

	} else if strings.Contains(args[1], "reset") {

		index = rand.Intn(len(quoteResetAudio))
		name = strings.Split(quoteResetAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteResetAudio[index])
		text = quoteResetText[index]

	} else if strings.Contains(args[1], "show") {

		index = rand.Intn(len(quoteShowAudio))
		name = strings.Split(quoteShowAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteShowAudio[index])
		text = quoteShowText[index]

	} else if strings.Contains(args[1], "merge") {

		index = rand.Intn(len(quoteMergeAudio))
		name = strings.Split(quoteMergeAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteMergeAudio[index])
		text = quoteMergeText[index]

	} else if strings.Contains(args[1], "checkout") {

		index = rand.Intn(len(quoteCheckoutAudio))
		name = strings.Split(quoteCheckoutAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteCheckoutAudio[index])
		text = quoteCheckoutText[index]

	} else if strings.Contains(args[1], "stash") {

		index = rand.Intn(len(quoteStashAudio))
		name = strings.Split(quoteStashAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quoteStashAudio[index])
		text = quoteStashText[index]

	} else if strings.Contains(args[1], "push") {

		index = rand.Intn(len(quotePushAudio))
		name = strings.Split(quotePushAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quotePushAudio[index])
		text = quotePushText[index]

	} else if strings.Contains(args[1], "pull") {

		index = rand.Intn(len(quotePullAudio))
		name = strings.Split(quotePullAudio[index], "/")
		url = fmt.Sprintf("%s%s", baseURl, quotePullAudio[index])
		text = quotePullText[index]

	} else {
		//os.Exit(1)
	}

	// locate the temp folder and build the full qualified name of the audio file
	temp := os.TempDir()
	audioFile := fmt.Sprintf("%s/%s", temp, name)

	// dowload the ogg file from the web
	downloadFromUrl(url, audioFile)

	runGit(args)

	// output text and audio
	if len(text) > 0 {
		fmt.Println("NARRATOR:", text)
		playAudio(audioFile)

		// remove ogg file
		os.Remove(audioFile)
	}

	return
}
