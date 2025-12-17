import { Component, signal } from '@angular/core';
import { Header } from './components/header/header';
import { Decoder } from './components/decoder/decoder';
import { Encoder } from './components/encoder/encoder';

@Component({
  selector: 'app-root',
  imports: [Header, Encoder, Decoder],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('geostego-web');
}
