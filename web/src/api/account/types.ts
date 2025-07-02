export type UpdatePreferencesRequest = {
  language?: string;
  theme?: string;
  timezone?: string;
};

export type RegisterRequest = {
  username: string;
  email: string;
  password: string;
  name: string;
};
