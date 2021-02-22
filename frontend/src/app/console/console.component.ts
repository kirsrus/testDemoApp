import {Component, OnInit} from '@angular/core';
import {DomSanitizer, SafeResourceUrl} from '@angular/platform-browser';

@Component({
  selector: 'app-console',
  templateUrl: './console.component.html',
  styleUrls: ['./console.component.scss']
})
export class ConsoleComponent implements OnInit {

  address: SafeResourceUrl;

  constructor(private dom: DomSanitizer) {
  }

  ngOnInit(): void {
    const url = new URL(window.location.href);
    url.port = '8080';
    url.pathname = '';
    const addr = url.toString();
    this.address = this.dom.bypassSecurityTrustResourceUrl(addr);
    console.log(this.address);
  }

}
