interface AlterFeature {
  properties: {
    event: string;
    areaDesc: string;
    severity: string;
    status: string;
    headline: string;
  };
}

interface ForecastPeriod {
  name: string;
  temperature: number;
  temperatureUnit: string;
  windSpeed: string;
  windDirection: string;
  shortForecast: string;
}

export type OptionalAlterFeature = Partial<AlterFeature>;
export type OptionalForecastPeriod = Partial<ForecastPeriod>;

export type AlertsResponse = Record<"features", OptionalAlterFeature[]>;
export type PointsResponse = Record<"properties", { forecast?: string }>;
export type ForecastResponse = Record<
  "properties",
  Record<"periods", OptionalForecastPeriod[]>
>;
