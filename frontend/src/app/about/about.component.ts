import {Component, OnInit} from '@angular/core';
import {SharedService} from '../services/shared.service';

@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent implements OnInit {
  title = {
    ru: 'О программе',
    en: 'About',
    zh: '关于程序 ',
  };

  constructor(private sharedService: SharedService) {
  }

  ngOnInit(): void {
    this.sharedService.setTitle(this.title);
  }

}
