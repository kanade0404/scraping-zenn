import {Firestore} from '@google-cloud/firestore';
import {CONFIG} from './config';

const {PROJECT_ID, GOOGLE_CREDENTIALS} = CONFIG;

type GetFirestore = () => Firestore;
export const getFirestore: GetFirestore = () => {
  return new Firestore({
    project_id: PROJECT_ID,
    keyFilename: GOOGLE_CREDENTIALS,
  });
};
