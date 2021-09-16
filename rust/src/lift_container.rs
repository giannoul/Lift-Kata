pub mod lift {


    // Directions ..
    pub const UP: i32 = 0;
    pub const DOWN: i32 = 1;

    // Call ..
    pub struct Call {
	    pub floor: i32,
	    pub direction: i32
    }

    // Lift ..
    pub struct Lift {
        pub id: String,
        pub floor: i32,
        pub requests:  Vec<i32>,
        pub doors_open: bool
    }

    impl Lift {
        pub fn print_lift(&self) -> String {
            let mut lift = format!("[{}{}]", self.has_request_on_floor(), self.id); // We set the doors closed by default
            if self.doors_open == true {
                lift = format!("]{}{}[", self.has_request_on_floor(), self.id);
            }
            return lift.to_string()
        }

        pub fn direction(&self) -> i32 {
            let mut d = UP;
            if self.floor > self.requests[0] {
                d = DOWN;
            }
            return d;
        }

        pub fn push_request(&mut self, req: i32) {
            self.requests.push(req);    
        }

        pub fn make_move(&mut self) {
            if self.requests.len() == 0 {
                return;
            }

            if self.direction() == UP {
                self.floor += 1
            }else {
                self.floor -= 1
            }
        }

        fn has_request_on_floor(&self) -> String {
            let mut r = "";
            for i in 0..self.requests.len() {
                if self.requests[i] == self.floor {
                    r = "*";
                    break;
                }
            }
            return r.to_string();
        }

    }

    // System ..
    pub struct System {
        pub floors: Vec<i32>,
        pub lifts:  Vec<Lift>,
        pub calls:  Vec<Call>
    }

    impl System {
                
        fn print_lifts_for_floor(&self, f: i32) -> String {
            let mut line = "".to_string();
            for i in 0..self.lifts.len() {
                let mut lift = "   ".to_string();
                if self.lifts[i].floor == f {
                    lift = self.lifts[i].print_lift();
                }else{
                    for j in 0..self.lifts[i].requests.len() { // Print lift requests for the current floor
                        if self.lifts[i].requests[j] == f {
                            lift = " * ".to_string();
                        }
                    }
                }
                line.push_str(&lift);
                line.push_str(&format!(" {}", f));
            }
            return line.to_string()
        }

        fn print_calls(&self, f: i32) -> String {
            let mut d = "";
            for i in 0..self.calls.len() {
                d = "   ";
                if self.calls[i].floor == f {
                    match self.calls[i].direction {
                        0 => d = " ^ ",
                        1 => d = " v ",
                        _ => d = "   "
                    }
                }
            }  
            return d.to_string()          
        }

        pub fn print_lifts(&self) {
            for i in  (0..self.floors.len()).rev() {
                println!("{0}{1}{2}", &self.floors[i],self.print_calls(i as i32),self.print_lifts_for_floor(i as i32));
            }
        }


        fn move_lifts(&mut self){
            for i in 0..self.lifts.len() {
                self.lifts[i].make_move();
            }
        }

        pub fn tick(&mut self){
            self.move_lifts();
            self.print_lifts();
        }


    }

}

