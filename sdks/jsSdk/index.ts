
import 'dotenv/config';

import { GetMeClient } from './src/client';



import { server } from "./src/server"

server.listen( process.env.PORT || 3000, () => {

  console.log(`Server is running on port ${process.env.PORT || 3000}`);

  

} );

export { GetMeClient }





