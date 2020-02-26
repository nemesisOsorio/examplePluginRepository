package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type PluginReleaseEvent struct {
    Org string
    Repo string
    Released Plugin
}

type Plugin struct {
    Id string           `json:"id"`
    Description string  `json:"description"`
    Provider string     `json:"provider"`
    Releases []Release  `json:"releases"`
}

type Release struct {
    Version string      `json:"version"`
    Date string         `json:"date"`
    Requires string     `json:"requires"`
    Sha512sum string    `json:"sha512sum"`
    State string        `json:"state"`
    Url string          `json:"url"`
}

func main() {

    pluginReleaseJson := []byte(os.Args[1])
    var pluginReleaseEvent PluginReleaseEvent
    pluginReleaseErr := json.Unmarshal(pluginReleaseJson, &pluginReleaseEvent)
    check(pluginReleaseErr)

    pluginsJson, pluginsJsonReadErr := ioutil.ReadFile("plugins.json")
    check(pluginsJsonReadErr)
    var plugins []Plugin
    pluginsErr := json.Unmarshal(pluginsJson, &plugins)
    check(pluginsErr)

    updatedPlugins := addReleaseToPlugins(pluginReleaseEvent, plugins)
    updatedPluginsJson, pluginReleaseErr := json.MarshalIndent(updatedPlugins, "  ", "  ")
    check(pluginReleaseErr)
    pluginsJsonWriteErr := ioutil.WriteFile("plugins.json", updatedPluginsJson, 0644)
    check(pluginsJsonWriteErr)
}


func addReleaseToPlugins(releaseEvent PluginReleaseEvent, existingPlugins []Plugin) []Plugin {
    releasedPlugin := releaseEvent.Released
    release := releasedPlugin.Releases[0]
    version := release.Version[1:]
    release.Url = "https://github.com/" + releaseEvent.Org + "/" + releaseEvent.Repo + "/releases/download/v" + version + "/" + releasedPlugin.Id + "-" + version + ".zip"
    releasedPlugin.Releases = []Release {
        release,
    }

    for ip, existingPlugin := range existingPlugins {
        if existingPlugin.Id == releasedPlugin.Id {
            for ir, existingRelease := range existingPlugin.Releases {
                if existingRelease.Version == release.Version {
                    existingPlugin.Releases = append(existingPlugin.Releases[:ir], existingPlugin.Releases[ir+1:]...)
                }
            }
            releasedPlugin.Releases = append(releasedPlugin.Releases, existingPlugin.Releases...)
            existingPlugins[ip] = releasedPlugin
            return existingPlugins
        }
    }

    return append(existingPlugins, releasedPlugin)
}


