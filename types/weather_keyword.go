package types

type WeatherKeyword string

const (
	Sunny   WeatherKeyword = "sunny"
	Cloudy  WeatherKeyword = "cloudy"
	Rainy   WeatherKeyword = "rainy"
	Snowy   WeatherKeyword = "snowy"
	Unknown WeatherKeyword = "unknown" // Default fallback
)

type WeatherCondition struct {
	Code           int            `json:"code"`
	Day            string         `json:"day"`
	Night          string         `json:"night"`
	Icon           int            `json:"icon"`
	WeatherKeyword WeatherKeyword `json:"weather_keyword"`
}

var WEATHER_CONDITIONS map[int]WeatherCondition = map[int]WeatherCondition{
	1000: {1000, "Sunny", "Clear", 113, Sunny},
	1003: {1003, "Partly cloudy", "Partly cloudy", 116, Cloudy},
	1006: {1006, "Cloudy", "Cloudy", 119, Cloudy},
	1009: {1009, "Overcast", "Overcast", 122, Cloudy},
	1030: {1030, "Mist", "Mist", 143, Cloudy},
	1063: {1063, "Patchy rain possible", "Patchy rain possible", 176, Rainy},
	1066: {1066, "Patchy snow possible", "Patchy snow possible", 179, Snowy},
	1069: {1069, "Patchy sleet possible", "Patchy sleet possible", 182, Snowy},
	1072: {1072, "Patchy freezing drizzle possible", "Patchy freezing drizzle possible", 185, Snowy},
	1087: {1087, "Thundery outbreaks possible", "Thundery outbreaks possible", 200, Rainy},
	1114: {1114, "Blowing snow", "Blowing snow", 227, Snowy},
	1117: {1117, "Blizzard", "Blizzard", 230, Snowy},
	1135: {1135, "Fog", "Fog", 248, Cloudy},
	1147: {1147, "Freezing fog", "Freezing fog", 260, Cloudy},
	1150: {1150, "Patchy light drizzle", "Patchy light drizzle", 263, Rainy},
	1153: {1153, "Light drizzle", "Light drizzle", 266, Rainy},
	1168: {1168, "Freezing drizzle", "Freezing drizzle", 281, Snowy},
	1171: {1171, "Heavy freezing drizzle", "Heavy freezing drizzle", 284, Snowy},
	1180: {1180, "Patchy light rain", "Patchy light rain", 293, Rainy},
	1183: {1183, "Light rain", "Light rain", 296, Rainy},
	1186: {1186, "Moderate rain at times", "Moderate rain at times", 299, Rainy},
	1189: {1189, "Moderate rain", "Moderate rain", 302, Rainy},
	1192: {1192, "Heavy rain at times", "Heavy rain at times", 305, Rainy},
	1195: {1195, "Heavy rain", "Heavy rain", 308, Rainy},
	1198: {1198, "Light freezing rain", "Light freezing rain", 311, Snowy},
	1201: {1201, "Moderate or heavy freezing rain", "Moderate or heavy freezing rain", 314, Snowy},
	1204: {1204, "Light sleet", "Light sleet", 317, Snowy},
	1207: {1207, "Moderate or heavy sleet", "Moderate or heavy sleet", 320, Snowy},
	1210: {1210, "Patchy light snow", "Patchy light snow", 323, Snowy},
	1213: {1213, "Light snow", "Light snow", 326, Snowy},
	1216: {1216, "Patchy moderate snow", "Patchy moderate snow", 329, Snowy},
	1219: {1219, "Moderate snow", "Moderate snow", 332, Snowy},
	1222: {1222, "Patchy heavy snow", "Patchy heavy snow", 335, Snowy},
	1225: {1225, "Heavy snow", "Heavy snow", 338, Snowy},
	1237: {1237, "Ice pellets", "Ice pellets", 350, Snowy},
	1240: {1240, "Light rain shower", "Light rain shower", 353, Rainy},
	1243: {1243, "Moderate or heavy rain shower", "Moderate or heavy rain shower", 356, Rainy},
	1246: {1246, "Torrential rain shower", "Torrential rain shower", 359, Rainy},
	1249: {1249, "Light sleet showers", "Light sleet showers", 362, Snowy},
	1252: {1252, "Moderate or heavy sleet showers", "Moderate or heavy sleet showers", 365, Snowy},
	1255: {1255, "Light snow showers", "Light snow showers", 368, Snowy},
	1258: {1258, "Moderate or heavy snow showers", "Moderate or heavy snow showers", 371, Snowy},
	1261: {1261, "Light showers of ice pellets", "Light showers of ice pellets", 374, Snowy},
	1264: {1264, "Moderate or heavy showers of ice pellets", "Moderate or heavy showers of ice pellets", 377, Snowy},
	1273: {1273, "Patchy light rain with thunder", "Patchy light rain with thunder", 386, Rainy},
	1276: {1276, "Moderate or heavy rain with thunder", "Moderate or heavy rain with thunder", 389, Rainy},
	1279: {1279, "Patchy light snow with thunder", "Patchy light snow with thunder", 392, Snowy},
	1282: {1282, "Moderate or heavy snow with thunder", "Moderate or heavy snow with thunder", 395, Snowy},
}
