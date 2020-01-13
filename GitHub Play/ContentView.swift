//
//  ContentView.swift
//  GitHub Play
//
//  Created by Academia on 1/13/20.
//  Copyright © 2020 Cate. All rights reserved.
//

import SwiftUI

struct ContentView: View {
    @State var showText = Color.white
    var body: some View {
        VStack {
            Image(systemName: "cloud")
                .resizable()
                .aspectRatio(contentMode: .fit)
                .padding(.horizontal, 100)
                .padding(.vertical, 10)
            Text("GitHub is available")
            Button(action:{
                self.showText = .black
            }) {
                Text("Learn more ...")
            }
            Text("Google 'Github Tutorial'").foregroundColor(showText)
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
