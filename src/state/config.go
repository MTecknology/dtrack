// Application Configuration
//
// ALL VARIABLES:
// - Define Runtime variables and types:	Application_Configuration
// - Default value of Runtime variables:	Load_Configuration
// - Map JSON Settings to Runtime variables:	Application_Configuration
// - Map Environment vars to Runtime variables:	Environment_Configation_Map
// - Help text showing JSON/EnvVar/Default:	Show_Help
package state

import (
	// DTrack
	"dtrack/log"

	// Standard
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Master object holding loaded configuration data
var Runtime Application_Configuration

// Map json configuration to Runtime
// Defaults set in load_config()
type Application_Configuration struct {
	Workspace              string   `json:"workspace"`
	Workspace_Keep_Temp    bool     `json:"workspace_keep_temp"`
	Record_Audio_Device    string   `json:"record_audio_device"`
	Record_Audio_Options   []string `json:"record_audio_options"`
	Record_Video_Device    string   `json:"record_video_device"`
	Record_Video_Options   []string `json:"record_video_options"`
	Record_Video_Advanced  []string `json:"record_video_advanced"`
	Record_Inspect_Models  []string `json:"record_inspect_models"`
	Has_Models             bool
	Record_Inspect_Backlog int     `json:"record_inspect_backlog"`
	Record_Duration        string  `json:"record_duration"`
	Train_Target           float64 `json:"train_target"`
	Train_Rate             float64 `json:"train_rate"`
	Train_Momentum         float64 `json:"train_momentum"`
	Train_Dropout          float64 `json:"train_dropout"`
}

// Map environment variables to Runtime
var Environment_Configation_Map = map[string]string{
	"DTRACK_WORKSPACE":              "Workspace",
	"DTRACK_WORKSPACE_KEEP_TEMP":    "Workspace_Keep_Temp",
	"DTRACK_RECORD_AUDIO_DEVICE":    "Record_Audio_Device",
	"DTRACK_RECORD_AUDIO_OPTIONS":   "Record_Audio_Options",
	"DTRACK_RECORD_VIDEO_DEVICE":    "Record_Video_Device",
	"DTRACK_RECORD_VIDEO_OPTIONS":   "Record_Video_Options",
	"DTRACK_RECORD_VIDEO_ADVANCED":  "Record_Video_Advanced",
	"DTRACK_RECORD_INSPECT_MODELS":  "Record_Inspect_Models",
	"DTRACK_RECORD_INSPECT_BACKLOG": "Record_Inspect_Backlog",
	"DTRACK_RECORD_DURATION":        "Record_Duration",
	"DTRACK_TRAIN_TARGET":           "Train_Target",
	"DTRACK_TRAIN_RATE":             "Train_Rate",
	"DTRACK_TRAIN_MOMENTUM":         "Train_Momentum",
	"DTRACK_TRAIN_DROUPOUT":         "Train_Dropout",
}

// Print information about configution file
func Show_Help() {
	fmt.Println("\nConfiguration Options:")
	fmt.Println("  Config.JSON Key\t\tEnvironment Variable\t\tDefault Value")
	fmt.Println("  ---------------\t\t--------------------\t\t-------------")
	fmt.Println("  workspace\t\t\tDTRACK_WORKSPACE\t\t_workspace")
	fmt.Println("  workspace_keep_temp\tDTRACK_WORKSPACE_KEEP_TEMP\tfalse")
	fmt.Println("  record_audio_device\t\tDTRACK_RECORD_AUDIO_DEVICE\tplughw")
	fmt.Println("  record_audio_options\t\tDTRACK_RECORD_AUDIO_OPTIONS\t[\"-f\", \"alsa\"]")
	fmt.Println("  record_video_device\t\tDTRACK_RECORD_VIDEO_DEVICE\t/dev/video0")
	fmt.Print("  record_video_options\t\tDTRACK_RECORD_VIDEO_OPTIONS")
	fmt.Println("\t[\"-f\", \"v4l2\", \"-framerate\", \"5\"]")
	fmt.Println("  record_video_advanced\t\tDTRACK_RECORD_VIDEO_ADVANCED\tSee Documentation")
	fmt.Println("  record_duration\t\tDTRACK_RECORD_DURATION\t\t00:10:00  (10 minutes)")
	fmt.Println("  record_inspect_models\t\tDTRACK_RECORD_INSPECT_MODELS\t[]")
	fmt.Println("  record_inspect_backlog\tDTRACK_RECORD_INSPECT_BACKLOG\t5")
	fmt.Println("  record_inspect_segment\tDTRACK_RECORD_INSPECT_SEGMENT\t-1")
	fmt.Println("  train_target\t\t\tDTRACK_TRAIN_TARGET\t\t0.95")
	fmt.Println("  train_rate\t\t\tDTRACK_TRAIN_RATE\t\t0.001")
	fmt.Println("  train_momentum\t\tDTRACK_TRAIN_MOMENTUM\t\t0.9")
	fmt.Println("  train_dropout\t\t\tDTRACK_TRAIN_DROUPOUT\t\t0.2")
	fmt.Println("\nOption Priority:")
	fmt.Println("  Defaults -> Config.json -> Environment")
}

// Loads Runtime configuration data into current state
func Load_Configuration(config_path string) {
	// Default configuration values
	cfg := Application_Configuration{
		Workspace:            "_workspace",
		Workspace_Keep_Temp:  false,
		Record_Audio_Device:  "plughw",
		Record_Audio_Options: []string{"-f", "alsa"},
		Record_Video_Device:  "/dev/video0",
		Record_Video_Options: []string{
			"-f", "v4l2", "-framerate", "15"},
		Record_Video_Advanced: []string{
			"-crf", "23", "-preset", "fast", "-maxrate", "3M", "-bufsize", "24M",
			"-tune", "zerolatency", "-filter_complex", "[1:v]drawtext" +
				"=fontfile=/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf" +
				":text=%{localtime}:fontcolor=red@0.9:x=7:y=7:fontsize=48[dtstamp]"},
		Record_Duration:        "00:10:00",
		Record_Inspect_Models:  []string{},
		Record_Inspect_Backlog: 5,
		Train_Target:           0.95,
		Train_Rate:             0.001,
		Train_Momentum:         0.9,
		Train_Dropout:          0.2,
	}

	// Check for configuration file
	log.Debug("Loading configuration values from: %s", config_path)
	if _, err := os.Stat(config_path); err != nil {
		log.Info("Configuration file not found; using defaults.")
		Runtime = cfg
		return
	}

	// Load configuration file
	file_data, err := os.ReadFile(config_path)
	if err != nil {
		log.Die("Error opening configuration file; ABORT!")
	}

	// Merge configuration values into cfg
	if err := json.Unmarshal(file_data, &cfg); err != nil {
		log.Die("Failed to parse configuration as JSON; ABORT!")
	}

	// Search for known environment variables
	for env_key, conf_field := range Environment_Configation_Map {
		if env_value := os.Getenv(env_key); env_value != "" {
			log.Debug("Environment variable found: %s", env_key)
			field := reflect.ValueOf(&cfg).Elem().FieldByName(conf_field)
			if !field.IsValid() || !field.CanSet() {
				log.Die("Invalid field: %s", conf_field)
			}

			// Merge environment variables into cfg
			switch field.Kind() {
			case reflect.String:
				field.SetString(env_value)
			case reflect.Int:
				if intVal, err := strconv.Atoi(env_value); err == nil {
					field.SetInt(int64(intVal))
				} else {
					log.Die("%s is not Integer", env_key)
				}
			case reflect.Float64:
				if intVal, err := strconv.Atoi(env_value); err == nil {
					field.SetFloat(float64(intVal))
				} else {
					log.Die("%s is not Float64", env_key)
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(env_value); err == nil {
					field.SetBool(boolVal)
				} else {
					log.Die("%s is not Boolean", env_key)
				}
			default:
				log.Die("Unexpected field type for %s", conf_field)
			}
		}
	}

	// Helper variables
	cfg.Has_Models = len(cfg.Record_Inspect_Models) > 0

	// Update session variables
	Runtime = cfg
}
