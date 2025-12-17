import { Component, signal } from '@angular/core';
import { StegoApiService } from '../../services/stego-api.service';

@Component({
  selector: 'app-decoder',
  imports: [],
  templateUrl: './decoder.html',
  styleUrl: './decoder.css',
})
export class Decoder {
  decodedMessage = signal('');
  errorMessage = signal('');
  isLoading = signal(false);
  selectedFile: File | null = null;

  constructor(private api: StegoApiService) {}

  onFileSelected(event: any) {
    this.selectedFile = event.target.files[0];
    this.errorMessage.set('');
    this.decodedMessage.set('');
  }

  onDecode() {
    if (!this.selectedFile) return;

    this.isLoading.set(true);
    this.errorMessage.set('');
    
    this.api.decode(this.selectedFile).subscribe({
      next: (res) => {
        this.decodedMessage.set(res.message);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.errorMessage.set(err.error?.error || 'Access Denied');
        this.isLoading.set(false);
      }
    });
  }
}
