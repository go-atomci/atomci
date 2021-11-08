import config from 'config';
import fetchAPI, { setServices } from 'za-fetch-api';

setServices(config.API);

export default fetchAPI;
