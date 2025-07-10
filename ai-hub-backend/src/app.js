import express from 'express';
import cors from 'cors';
import connectDB from './services/mongoDB/connectDB.js';
import rootRoutes from './routes/rootRoutes.js';
import loadEnv from './config/loadEnv.js';

const app = express();

connectDB();             

app.use(cors());
app.use(express.json());
app.use('/api', rootRoutes);
 
const PORT = loadEnv.PORT;
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
