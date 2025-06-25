export type LoginRequest = {
  login: string;
  password: string;
};

export type RegisterRequest = {
  username: string;
  email: string;
  password: string;
  name: string;
};

export type UpdatePreferencesRequest = {
  language?: string;
  theme?: string;
  timezone?: string;
};
