import type {HttpFunction} from '@google-cloud/functions-framework/build/src/functions';
import {start} from './scraping';

export const main: HttpFunction = async (req, res) => {
  console.log('start');
  const {name} = req.query;
  try {
    console.log(`before scraping name:${name}`);
    if (!name) {
      res.status(400).send({message: "Query strings 'name' is required."})
      return
    }
    const articles = await start(name.toString());
    console.log(`success scraping name:${name}`);
    res.send({message: 'success', data: articles});
  } catch (e) {
    console.error(e);
    res.status(400).send({message: 'error occurred.', error: e});
  }
};
