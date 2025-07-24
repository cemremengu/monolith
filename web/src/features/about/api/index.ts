export type VersionInfo = {
  version: string;
  commit: string;
  dateBuilt: string;
};

export const getVersionInfo = async (): Promise<VersionInfo> => {
  const response = await fetch("/api/version");
  if (!response.ok) {
    throw new Error("Failed to fetch version info");
  }
  return response.json();
};
