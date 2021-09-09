mod lift_container;

fn main() {
    let f: Vec<i32> = [0,1,2,3,4,5,6].to_vec();
    let mut c = Vec::with_capacity(20);
    c.push(lift_container::lift::Call{
        floor: 4,
        direction: lift_container::lift::DOWN  
    });

    c.push(lift_container::lift::Call{
        floor: 1,
        direction: lift_container::lift::UP  
    });

    let mut l = Vec::with_capacity(20);
    l.push(    
        lift_container::lift::Lift{
            id: "A".to_string(),
            floor: 2,
            requests:  [].to_vec(),
            doors_open: false
    });

    let mut system = lift_container::lift::System{
        calls: c,
        floors: f,
        lifts: l
    };
    
    //system.print_all();
    system.print_lifts();
    system.lifts[0].push_request(3);
    
    system.lifts[0].push_request(1);
    system.lifts[0].push_request(5);

    println!("");
    
    for i in 0..system.lifts[0].requests.len() {
        println!("{0}",system.lifts[0].requests[i]);
    }
}
