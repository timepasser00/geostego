import { Component, signal } from '@angular/core';
import { StegoApiService } from '../../services/stego-api.service';
import { FileDownloadService } from '../../services/file-download.service';

@Component({
  selector: 'app-encoder',
  standalone: true,
  imports: [],
  templateUrl: './encoder.html',
  styleUrl: './encoder.css',
})
export class Encoder {
  message = signal<string>('')
  isLoading = signal(false);
  selectedFile: File | null = null;

  constructor(private api: StegoApiService, private fileService: FileDownloadService) {}

  onFileSelected(event: any) {
    this.selectedFile = event.target.files[0];
  }

  onMessageChange(event : Event) {

    const inputValue = (event.target as HTMLTextAreaElement).value
    this.message.set(inputValue);
  }

  onEncode() {
    if (!this.selectedFile || !this.message()) return;
    
    this.isLoading.set(true);
    this.api.encode(this.selectedFile, this.message()).subscribe({
      next: (response) => {
        this.fileService.downloadFile(response, 'secret_stego.png');
        this.isLoading.set(false);
        this.message.set('');
      },
      error: () => {
        alert('Encoding Failed');
        this.isLoading.set(false);
      }
    });
  }
}
