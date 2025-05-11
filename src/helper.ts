import { OptionalAlterFeature } from "./interface.js";

export const makeNWSRequest = async <T>(
  userAgent: string,
  url: string
): Promise<T | null> => {
  const headers = {
    "User-Agent": userAgent,
    Accept: "application/geo+json",
  };

  try {
    const response = await fetch(url, { headers });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return data as T;
  } catch (error) {
    console.error(error);
    return null;
  }
};

export const formatAlert = (feature: OptionalAlterFeature): string => {
  const props = feature.properties;
  return [
    `Event : ${getPropsReturnValue(props?.event)}`,
    `Area : ${getPropsReturnValue(props?.areaDesc)}`,
    `Severity : ${getPropsReturnValue(props?.severity)}`,
    `Status : ${getPropsReturnValue(props?.status)}`,
    `Headline : ${getPropsReturnValue(props?.headline, "No headline")}`,
    "---",
  ].join("\n");
};

const getPropsReturnValue = (target?: string, defaultValue = "Unknown") =>
  target || defaultValue;
