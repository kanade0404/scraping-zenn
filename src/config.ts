import {config} from 'dotenv';

config();

const {PROJECT_ID, GOOGLE_CREDENTIALS} = process.env;

export const CONFIG = {
  PROJECT_ID,
  GOOGLE_CREDENTIALS,
} as const;
