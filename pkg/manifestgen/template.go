package manifestgen

const manifestTemplate = `fx_version 'cerulean'
game 'gta5'

files {
    {{ if .HasCarcols -}}
    'data/carcols/*.meta',
    {{ end -}}
    {{ if .HasCarvariations -}}
    'data/carvariations/*.meta',
    {{ end -}}
    {{ if .HasContentUnlocks -}}
    'data/contentunlocks/*.meta',
    {{ end -}}
    {{ if .HasHandling -}}
    'data/handling/*.meta',
    {{ end -}}
    {{ if .HasVehicleLayouts -}}
    'data/vehiclelayouts/*.meta',
    {{ end -}}
    {{ if .HasVehicleModelsets -}}
    'data/vehiclemodelsets/*.meta',
    {{ end -}}
    {{ if .HasVehicles -}}
    'data/vehicles/*.meta',
    {{ end -}}
    {{ if .HasWeaponsFile -}}
    'data/weaponsfile/*.meta',
    {{ end -}}
    {{ if .HasAudio -}}
    {{ range $name := .AudioConfigs -}}
    'audioconfig/{{ $name }}_game.dat151.rel',
    'audioconfig/{{ $name }}_game.dat151.nametable',
    'audioconfig/{{ $name }}_sounds.dat54.rel',
    'audioconfig/{{ $name }}_sounds.dat54.nametable',
    {{ end -}}
    {{ range $pack := .AudioWavePacks -}}
    'sfx/dlc_{{ $pack }}/*.awc',
    'sfx/dlc_{{ $pack }}/*.awc2',
    {{ end -}}
    {{ end }}
}

{{ if .HasCarcols -}}
data_file 'CARCOLS_FILE' 'data/carcols/*.meta'
{{ end -}}
{{ if .HasCarvariations -}}
data_file 'VEHICLE_VARIATION_FILE' 'data/vehiclevariations/*.meta'
{{ end -}}
{{ if .HasContentUnlocks -}}
data_file 'CONTENT_UNLOCKING_META_FILE' 'data/contentunlocks/*.meta'
{{ end -}}
{{ if .HasHandling -}}
data_file 'HANDLING_FILE' 'data/handling/*.meta'
{{ end -}}
{{ if .HasVehicleLayouts -}}
data_file 'VEHICLE_LAYOUTS_FILE' 'data/vehiclelayouts/*.meta'
{{ end -}}
{{ if .HasVehicleModelsets -}}
data_file 'AMBIENT_VEHICLE_MODEL_SET_FILE' 'data/vehiclemodelsets/*.meta'
{{ end -}}
{{ if .HasVehicles -}}
data_file 'VEHICLE_METADATA_FILE' 'data/vehicles/*.meta'
{{ end -}}
{{ if .HasWeaponsFile -}}
data_file 'WEAPONINFO_FILE' 'data/weaponsfile/*.meta'
{{ end -}}
{{ if .HasAudio -}}
{{ range $name := .AudioConfigs -}}
data_file 'AUDIO_GAMEDATA' 'audioconfig/{{ $name }}_game.dat'
data_file 'AUDIO_SOUNDDATA' 'audioconfig/{{ $name }}_sounds.dat'
{{ end -}}
{{ range $pack := .AudioWavePacks -}}
data_file 'AUDIO_WAVEPACK' 'sfx/dlc_{{ $pack }}'
{{ end -}}
{{ end }}`
