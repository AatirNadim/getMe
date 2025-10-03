import express from 'express';


import 'dotenv/config';
import { GetMeClient } from './client';






const app = express();


const port = process.env.PORT || 3000;

const client = new GetMeClient();


app.get('/get', async (req, res) => {

  console.log('Received request', req.query);

  const key = req.query.key as string;

  if (!key) {

    return res.status(400).send({ error: 'Key is required' });

  }
  

  console.log('Fetching value for key:', key);
  const value = await client.get(key);
  
  if (value === null) {
    return res.status(404).send({ error: 'Key not found' });
  }

  res.send(value);

})


app.post('/put', async (req, res) => {
  console.log('Received request', req.query);

  const key = req.query.key as string;
  const value = req.query.value as string;

  console.log('Storing key-value pair:', key, value);

  if (!key || !value) {
    return res.status(400).send({ error: 'Key and value are required' });
  }

  await client.put(key, value);

  res.send({ success: true });
})


app.delete('/delete', async (req, res) => {
  console.log('Received request', req.query);

  const key = req.query.key as string;

  console.log('Deleting key:', key);

  if (!key) {
    return res.status(400).send({ error: 'Key is required' });
  }

  await client.delete(key);

  res.send({ success: true });
})


export { app as server }