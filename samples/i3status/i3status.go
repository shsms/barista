// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// i3status is a port of the default i3status configuration to barista.
package main

import (
	"time"

	"barista.run"
	"barista.run/bar"
	"barista.run/colors"
	"barista.run/format"
	"barista.run/modules/battery"
	"barista.run/modules/clock"
	"barista.run/modules/diskspace"
	"barista.run/modules/sysinfo"
	"barista.run/modules/volume"
	"barista.run/modules/volume/pulseaudio"
	"barista.run/outputs"
)

func main() {
	colors.LoadFromMap(map[string]string{
		"good":     "#50FA7B",
		"bad":      "#FF5555",
		"degraded": "#F1FA8C",
	})

	barista.Add(diskspace.New("/").Output(func(i diskspace.Info) bar.Output {
		out := outputs.Textf(" 💾 %s ", format.IBytesize(i.Available))
		switch {
		case i.AvailFrac() < 0.2:
			out.Color(colors.Scheme("bad"))
		case i.AvailFrac() < 0.33:
			out.Color(colors.Scheme("degraded"))
		}
		return out
	}))

	statusName := map[battery.Status]string{
		battery.Charging:    "CHR",
		battery.Discharging: "BAT",
		battery.NotCharging: "NOT",
		battery.Unknown:     "UNK",
	}
	barista.Add(battery.All().Output(func(b battery.Info) bar.Output {
		if b.Status == battery.Disconnected {
			return nil
		}
		if b.Status == battery.Full {
			return outputs.Text(" 🔋 FULL ")
		}
		out := outputs.Textf(" 🔋 %s %d%% %s ",
			statusName[b.Status],
			b.RemainingPct(),
			b.RemainingTime())
		if b.Discharging() {
			if b.RemainingPct() < 20 || b.RemainingTime() < 30*time.Minute {
				out.Color(colors.Scheme("bad"))
			}
		}
		return out
	}))

	barista.Add(sysinfo.New().Output(func(i sysinfo.Info) bar.Output {
		out := outputs.Textf(" 📈 %.2f ", i.Loads[0])
		if i.Loads[0] > 3.0 {
			out.Color(colors.Scheme("bad"))
		} else if i.Loads[0] > 2.0 {
			out.Color(colors.Scheme("degraded"))
		}
		return out
	}))

	barista.Add(volume.New(pulseaudio.DefaultSink()).Output(func(v volume.Volume) bar.Output {
		if v.Mute {
			return outputs.Text(" 🔇 ")
		}
		icon := "🔉"
		if v.Pct() > 60 {
			icon = "🔊"
		}
		return outputs.Textf(" %s %03d ", icon, v.Pct())
	}))

	barista.Add(clock.Local().OutputFormat(" ⌚ 2006-01-02 15:04:05 "))

	panic(barista.Run())
}
