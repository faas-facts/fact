package fact;

import java.util.HashMap;
import java.util.Map;
import com.amazonaws.services.lambda.runtime.Context;
import io.github.fact.Fact;
import io.github.fact.FactConfiguration;
import io.github.fact.FactConfigurationBuilder;
import io.github.fact.io.TCPSender;

import static java.lang.Thread.sleep;


public class TcpFactTest {
    static {
        FactConfigurationBuilder builder= new FactConfigurationBuilder();
        TCPSender sender = new TCPSender();
        FactConfiguration conf= builder.setIo(sender).setSendOnUpdate(true).createFactConfiguration();
        Fact.boot(conf);
        System.out.println("finished building fact conf");

    }

    public String handleRequest(Map<String, String> event, Context context) {
        Fact.start(context,event);
        System.out.println("starting completed");
        try {
            sleep(1000);
        } catch (InterruptedException e) { }
        System.out.println("first nap is done");
        Map<String,String>tags= new HashMap<String, String>();
        tags.put("experiment","#1");
        Fact.update(context,"update testing",tags);
        System.out.println("update complete");
        try {
            sleep(1000);
        } catch (InterruptedException e) { }
        System.out.println("second nap is done");
        Fact.done(context,"im done with this","no more args");
        System.out.println("all good and done");
        return null;
    }
}
