package fact;

import java.util.HashMap;
import java.util.Map;
import com.amazonaws.services.lambda.runtime.Context;
import io.github.fact.Fact;
import io.github.fact.FactConfiguration;
import io.github.fact.FactConfigurationBuilder;

import static java.lang.Thread.sleep;


public class ConsoleFactTest {

    public String handleRequest(Map<String, String> event, Context context) {
        FactConfigurationBuilder builder= new FactConfigurationBuilder();
        FactConfiguration conf= builder.createLazyLogger();
        Fact.boot(conf);
        Fact.start(context,event);
        try {
            sleep(1000);
        } catch (InterruptedException e) { }
        Map<String,String>tags= new HashMap<String, String>();
        tags.put("experiment","#1");
        Fact.update(context,"update testing",tags);
        try {
            sleep(1000);
        } catch (InterruptedException e) { }
        Fact.done(context,"im done with this","no more args");
        return null;
    }
}
