import { HttpClient, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environments';

@Injectable({
  providedIn: 'root',
})

export class StegoApiService {
  private baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  encode(image: File, message: string): Observable<HttpResponse<Blob>> {
    const formData = new FormData();
    formData.append('image', image);
    formData.append('message', message);

    return this.http.post(`${this.baseUrl}/encode`, formData, {
      responseType: 'blob',
      observe: 'response'
    });
  }

  decode(image: File): Observable<any> {
    const formData = new FormData();
    formData.append('image', image);

    return this.http.post(`${this.baseUrl}/decode`, formData);
  }
}
