import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import {Apollo, gql} from 'apollo-angular';
import {GraphQLError} from 'graphql';
import {MatTable, MatTableDataSource} from '@angular/material/table';
import {MatPaginator} from '@angular/material/paginator';
import {SharedService} from '../services/shared.service';

// Информация о хосте (GraphQL)
export interface HostInfoQuery {
  hostInfo: {
    uptime: number
    hostName: string
    hostID: string
    os: string
    platform: string
    platformFamily: string
    platformVersion: string
    kernelVersion: string
    kernelArch: string
  };
}

// Информация о процессах на хосте (GraphQL)
export interface ProcessInfoQuery {
  processInfo: {
    name: string
    pid: number
    vms: number
    hwm: number
  }[];
}

// Элементы дааннхы таблицы с отображением процессвов на хсоте
export interface TableProcess {
  name: string;
  pid: number;
  vms: number;
  hwm: number;
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit, AfterViewInit {
  title = {
    ru: 'МЛТ - медицинские технологии - панель управления',
    en: 'MLT - medical technologies - control panel',
    zh: 'MLT - 医疗技术-控制面板',
  };


  // Информация о хосте
  hostInfo: HostInfoQuery = {
    hostInfo: {
      uptime: 0,
      hostName: '',
      hostID: '',
      os: '',
      platform: '',
      platformFamily: '',
      platformVersion: '',
      kernelVersion: '',
      kernelArch: '',
    }
  };

  // Информация о процессах хоста
  processInfo: ProcessInfoQuery;

  // Таблица с информацией о процессых на хосте
  tableProcessColumns: string[] = ['name', 'pid', 'vms', 'hwm'];
  tableProcessDataSource = new MatTableDataSource<TableProcess>([]);
  // Ссылка на таблицу
  @ViewChild(MatTable, {static: true}) table: MatTable<any>;

  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(
    private apollo: Apollo,
    private sharedService: SharedService
  ) {
  }

  ngOnInit(): void {
    this.sharedService.setTitle(this.title);
    this.loadHostInfo();
    this.loadProcessInfo();
  }

  ngAfterViewInit(): void {
    this.tableProcessDataSource.paginator = this.paginator;
  }

  loadHostInfo(): void {
    this.apollo.query<HostInfoQuery>({
      query: gql`
        query {
          hostInfo{
            uptime
            hostName
            hostID
            os
            platform
            platformFamily
            platformVersion
            kernelVersion
            kernelArch
            virtualizationSystem
            virtualizationRole
          }
        }
      `,
      fetchPolicy: 'no-cache',
    }).subscribe(
      ({data}) => {
        if (data) {
          this.hostInfo = data;
        }
        // console.log('HostInfoQuery:', data);
      },
      (err: GraphQLError) => {
        console.log('GraphQLError:', err);
      }
    );
  }

  loadProcessInfo(): void {
    this.apollo.query<ProcessInfoQuery>({
      query: gql`
        query {
          processInfo{
            name
            pid
            vms
          }
        }
      `,
      fetchPolicy: 'no-cache',
    }).subscribe(
      ({data}) => {
        if (data) {
          this.processInfo = data;

          // Обновляем таблицу данынми
          const tData: TableProcess[] = [];
          this.processInfo.processInfo.forEach(v => {
            tData.push({
              name: v.name,
              pid: v.pid,
              vms: v.vms,
              hwm: v.hwm,
            });
          });
          tData.sort((a, b) => a.name.localeCompare(b.name));
          this.tableProcessDataSource.data = tData;
          this.table.renderRows();
        }
      },
      (err: GraphQLError) => {
        console.log('GraphQLError:', err);
      }
    );
  }

}
