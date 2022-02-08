import {Article} from './type';
import {getFirestore} from './firestore';
import {FieldValue} from '@google-cloud/firestore';

type Save = (name: string, articles: Article[]) => Promise<void>;
export const save: Save = async (name, articles) => {
  const firestore = getFirestore();
  const document = firestore.collection(name).doc('zenn');
  await document.set({
    articles,
    timestamp: FieldValue.serverTimestamp(),
  });
};
