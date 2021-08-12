import http from 'k6/http';
import { sleep } from 'k6';
export let options = {
  vus: 100,
  duration: '30m',
};
export default function () {
  http.get('http://host.docker.internal:8000/blocknumber');
  //sleep(1);
}

