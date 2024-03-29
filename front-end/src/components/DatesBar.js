import { useEffect, useState } from "react";
import CalendarMenu from "./CalenderMenu";

const DAYSOFWEEK = {
    "1": "Monday",
    "2": "Tuesday",
    "3": "Wednesday",
    "4": "Thursday",
    "5": "Friday",
    "6": "Saturday",
    "0": "Sunday",
};

function DatesBar(props) {
    const [currentDate, setCurrentDate] = useState({});
    const [dates, setDates] = useState([]);
    const [selectedDate, setSelectedDate] = useState(null);
    const [selectedDateDoses, setselectedDateDoses] = useState([]);

    function setUpCurrent() {
        let cDate = new Date();

        let day = cDate.getDay();
        let date = cDate.getDate();
        let month = cDate.getMonth();

        if (Object.keys(currentDate).length === 0) {
            setCurrentDate({
                Day: DAYSOFWEEK[day],
                DayNum: date,
                Month: month + 1,
            });

            setDates([
                { Day: DAYSOFWEEK[day], DayNum: date, Month: month + 1 },
            ]);

            for (let i = 1; i < 4; i++) {
                let cDate = new Date();
                cDate.setDate(cDate.getDate() + i);

                let day = cDate.getDay();
                let date = cDate.getDate();
                let month = cDate.getMonth();

                setDates((prevState) => [
                    ...prevState,
                    { Day: DAYSOFWEEK[day], DayNum: date, Month: month + 1 },
                ]);
            }

            for (let i = 1; i < 4; i++) {
                let cDate = new Date();
                cDate.setDate(cDate.getDate() - i);

                let day = cDate.getDay();
                let date = cDate.getDate();
                let month = cDate.getMonth();

                setDates((prevState) => [
                    { Day: DAYSOFWEEK[day], DayNum: date, Month: month + 1 },
                    ...prevState,
                ]);
            }
        }
    }

    useEffect(() => {
        let timer = setInterval(() => setCurrentDate(new Date()), 1000);
        setUpCurrent();

        return function cleanup() {
            clearInterval(timer);
        };
    }, []);

    const prevDate = () => {

    };

    const nextDate = () => {
    
    };

    const handleMenu = (e, index) => {
        const [dayNum, month, selected] = e.currentTarget.id.split(' ');
        console.log(dayNum, month, selected);
        setSelectedDate({ dayNum, month });
        setselectedDateDoses([])
        props.schedule.forEach(element => {
            const date = new Date(element.datetime);
            if (date.getDate() === parseInt(dayNum) && date.getMonth() + 1 === parseInt(month)) {
                setselectedDateDoses((prevState) => [...prevState, element]);
            }
        });
    }

    if (dates.length === 0) {
        return <div>Loading</div>;
    }

    return (
        <div className="flex flex-col items-center justify-center space-y-2">
            {props.schedule == undefined ? <h1 class="text-white text-5xl text-center">You have no doses set, you may add some in the <span class="text-[#dddd55]">Doses</span> page</h1> :
            <div className="flex items-center space-x-2">
                <button
                    className="bg-slate-500 rounded-sm p-1"
                    onClick={null}
                >
                    Prev Date
                </button>

                {
                    dates.map((date) => {
                        let selected = false;

                        if (selectedDate !== null) {
                            if (date.DayNum === parseInt(selectedDate.dayNum) && date.Month === parseInt(selectedDate.month)) {
                                selected = true;
                            } else {
                                selected = false;
                            }
                        }

                        return(
                        <div key={date.DayNum} id={date.DayNum + ' ' + date.Month + ' ' + selected} className="flex-1 cursor-pointer" onClick={handleMenu}>
                            <div className={` ${selected ? 'bg-green-400' : 'bg-gray-200'} p-2 rounded-md text-center`}>
                                <p className="text-lg font-semibold">{date.Day}</p>
                                <p className="text-sm">
                                    {date.DayNum}/{date.Month}
                                </p>
                            </div>
                        </div>
                        )
                    })

                }

                <button
                    className="bg-slate-500 rounded-sm p-1"
                    onClick={nextDate}
                >
                    Next Date
                </button>
            </div>
        }

            {selectedDate !== null && props.schedule != undefined && <CalendarMenu schedule={selectedDateDoses} />}
        </div>


    );
}

export default DatesBar;
