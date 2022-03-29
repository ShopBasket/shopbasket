import { Injectable } from '@angular/core';
import { Inventory } from './inventory';
import { Observable } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { catchError, map, tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class InventoryService {
  private inventoryUrl = 'api/inventory'; 
  constructor(private http: HttpClient) { }
  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  };
  getInventory(): Observable<Inventory[]> {
    return this.http.get<Inventory[]>(this.inventoryUrl);
  }

  createInventory(name: String, description: String, price: Number): Observable<Inventory> {
    var inventory= {name:name, description:description, price:price, status:true}
    return this.http.post<Inventory>(this.inventoryUrl, inventory, this.httpOptions);
  }

}
