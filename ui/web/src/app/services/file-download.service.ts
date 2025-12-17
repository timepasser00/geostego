import { Injectable } from '@angular/core';
import { HttpResponse } from '@angular/common/http';

@Injectable({ providedIn: 'root' })
export class FileDownloadService {

  downloadFile(response: HttpResponse<Blob>, defaultName: string = 'file'): void {
    const blob = response.body;
    if (!blob) return;

    const filename = this.getFileNameFromHeaders(response.headers) || defaultName;

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    window.URL.revokeObjectURL(url);
  }

  private getFileNameFromHeaders(headers: any): string | null {
    const contentDisposition = headers.get('Content-Disposition');
    if (!contentDisposition) return null;

    const matches = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(contentDisposition);
    return (matches != null && matches[1]) ? matches[1].replace(/['"]/g, '') : null;
  }
}